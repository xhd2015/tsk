# Scenario

**Feature**: isolated TSK_HOME per test; build tsk binary once per process

```
# temp WorkRoot + {WorkRoot}/.tsk per leaf; process-local tsk binary (mutex memo)
# tskEnv strips CODEX_THREAD_ID / PI_CODING_AGENT / TSK_STATUS_FORMAT from parent env
# so bare `status` stays diagram-stable under host agents; leaves re-inject via ExtraEnv
tsk <subcommand> -> stdout/stderr + filesystem side effects under TSK_HOME
```

## Preconditions

- The tsk Go module root is two levels above the test tree (`github.com/xhd2015/tsk`).
- Go toolchain is available on PATH.
- Process-local binary: `getTskBin` builds once under an in-memory mutex into
  `os.MkdirTemp("", "tsk-doctest-bin-")` (not session disk flock).
- Child env always clears host-agent detection and `TSK_STATUS_FORMAT` unless a leaf sets `Request.ExtraEnv`.

## Context

Each leaf runs `tsk` with `TSK_HOME={WorkRoot}/.tsk` and `TSK_DATE=2026-07-09`. Helper `runTskCmd` / `Run` use `tskEnv(req)`, which strips `CODEX_THREAD_ID`, `PI_CODING_AGENT`, and `TSK_STATUS_FORMAT` from the parent process environment so diagram leaves do not flaky-auto to agent when CI runs under Codex/Grok. Auto-format leaves re-inject those vars via `req.ExtraEnv` (`KEY=value`). Per-leaf state is never shared across leaves.

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	"github.com/xhd2015/doctest/assert"
	"github.com/xhd2015/doctest/session"
)

const tskDate = "2026-07-09"

// Process-local tsk binary (one-process suite; in-memory mutex, not session flock).
var (
	tskBinMu   sync.Mutex
	tskBinPath string
	tskBinErr  error
	// tskModRoot set from d.DOCTEST_ROOT in root Setup.
	tskModRoot string
)

type taskJSON struct {
	ID           int              `json:"id"`
	Title        string           `json:"title"`
	Slug         string           `json:"slug"`
	Labels       []string         `json:"labels"`
	TopicPath    json.RawMessage  `json:"topic_path"`
	Stage        string           `json:"stage"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
	StageHistory []map[string]any `json:"stage_history"`
}

type clarifyBatch struct {
	BatchID string `json:"batch_id"`
	Status  string `json:"status"`
	Items   []struct {
		ID       string `json:"id"`
		Question string `json:"question"`
		Status   string `json:"status"`
	} `json:"items"`
}

type eventLine struct {
	TS      string   `json:"ts"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	ExitCode int     `json:"exit_code"`
}

func findModuleRoot(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func getTskBin(t *testing.T) string {
	t.Helper()
	tskBinMu.Lock()
	defer tskBinMu.Unlock()
	if tskBinPath != "" || tskBinErr != nil {
		if tskBinErr != nil {
			t.Fatal(tskBinErr)
		}
		return tskBinPath
	}
	if tskModRoot == "" {
		t.Fatal("tskModRoot unset; root Setup must run first")
	}
	dir, err := os.MkdirTemp("", "tsk-doctest-bin-")
	if err != nil {
		tskBinErr = err
		t.Fatal(err)
	}
	bin := filepath.Join(dir, "tsk")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/tsk")
	cmd.Dir = tskModRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		tskBinErr = fmt.Errorf("build tsk: %v\n%s", err, out)
		t.Fatal(tskBinErr)
	}
	tskBinPath = bin
	return bin
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if root := findModuleRoot(d.DOCTEST_ROOT); root != "" {
		tskModRoot = root
	} else {
		tskModRoot = filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	}
	workRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		return fmt.Errorf("resolve work root: %w", err)
	}
	req.WorkRoot = workRoot
	req.TskHome = filepath.Join(req.WorkRoot, ".tsk")
	ensureHelpersUsed()
	return os.MkdirAll(req.TskHome, 0o755)
}

// tskEnvBaseDrop are env keys always stripped from the parent process env so
// status auto-format does not flaky-select agent when CI runs under Codex/Grok/PI.
// Leaves re-inject via Request.ExtraEnv (e.g. CODEX_THREAD_ID=t1).
var tskEnvBaseDrop = map[string]struct{}{
	"CODEX_THREAD_ID":   {},
	"PI_CODING_AGENT":   {},
	"TSK_STATUS_FORMAT": {},
	"TSK_HOME":          {},
	"TSK_DATE":          {},
}

func envKey(entry string) string {
	if i := strings.IndexByte(entry, '='); i >= 0 {
		return entry[:i]
	}
	return entry
}

func filterEnvKeys(env []string, drop map[string]struct{}) []string {
	out := make([]string, 0, len(env))
	for _, e := range env {
		if _, ok := drop[envKey(e)]; ok {
			continue
		}
		out = append(out, e)
	}
	return out
}

// tskEnv builds the child process environment: parent env minus host-agent and
// format-override vars, plus isolated TSK_HOME/TSK_DATE, plus req.ExtraEnv.
// ExtraEnv keys also replace any remaining parent values with the same name.
func tskEnv(req *Request) []string {
	drop := make(map[string]struct{}, len(tskEnvBaseDrop)+len(req.ExtraEnv))
	for k := range tskEnvBaseDrop {
		drop[k] = struct{}{}
	}
	for _, e := range req.ExtraEnv {
		drop[envKey(e)] = struct{}{}
	}
	env := filterEnvKeys(os.Environ(), drop)
	env = append(env,
		"TSK_HOME="+req.TskHome,
		"TSK_DATE="+tskDate,
	)
	env = append(env, req.ExtraEnv...)
	return env
}

func captureCommandOutput(cmd *exec.Cmd) (*Response, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			return nil, err
		}
	}
	return &Response{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

func runTskCmd(t *testing.T, req *Request, args ...string) *Response {
	t.Helper()
	bin := getTskBin(t)
	cmd := exec.Command(bin, args...)
	cmd.Dir = req.WorkRoot
	cmd.Env = tskEnv(req)
	resp, err := captureCommandOutput(cmd)
	if err != nil {
		t.Fatalf("run tsk %v: %v", args, err)
	}
	return resp
}

func runTskOK(t *testing.T, req *Request, args ...string) *Response {
	t.Helper()
	resp := runTskCmd(t, req, args...)
	if resp.ExitCode != 0 {
		t.Fatalf("tsk %v: exit %d stderr=%q", args, resp.ExitCode, resp.Stderr)
	}
	return resp
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	s = b.String()
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	runes := []rune(s)
	if len(runes) > 64 {
		s = string(runes[:64])
		s = strings.Trim(s, "-")
	}
	return s
}

func taskDirName(id int, stage, title string) string {
	return fmt.Sprintf("%d-%s-%s", id, stage, slugify(title))
}

func inboxTaskRel(id int, stage, title string) string {
	return filepath.ToSlash(filepath.Join("inbox", taskDirName(id, stage, title)))
}

func topicTaskRel(topic string, id int, stage, title string) string {
	parts := strings.Split(topic, "/")
	all := append(parts, taskDirName(id, stage, title))
	return filepath.ToSlash(filepath.Join(append([]string{"topics"}, all...)...))
}

func taskAbs(req *Request, rel string) string {
	return filepath.Join(req.TskHome, filepath.FromSlash(rel))
}

func readIndex(t *testing.T, req *Request, id int) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(req.TskHome, "index", fmt.Sprintf("%d", id)))
	if err != nil {
		t.Fatalf("read index/%d: %v", id, err)
	}
	return strings.TrimSpace(string(data))
}

func readTaskJSON(t *testing.T, taskDir string) taskJSON {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(taskDir, "task.json"))
	if err != nil {
		t.Fatalf("read %s/task.json: %v", taskDir, err)
	}
	var task taskJSON
	if err := json.Unmarshal(data, &task); err != nil {
		t.Fatalf("parse task.json: %v", err)
	}
	return task
}

func readClarifyBatch(t *testing.T, taskDir string) clarifyBatch {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(taskDir, "clarify", "batch.json"))
	if err != nil {
		t.Fatalf("read clarify/batch.json: %v", err)
	}
	var batch clarifyBatch
	if err := json.Unmarshal(data, &batch); err != nil {
		t.Fatalf("parse batch.json: %v", err)
	}
	return batch
}

func parseCreatedAt(t *testing.T, raw string) time.Time {
	t.Helper()
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if ts, err := time.Parse(f, raw); err == nil {
			return ts
		}
	}
	t.Fatalf("parse created_at %q", raw)
	return time.Time{}
}

func createArgs(title string, topic string, labels []string) []string {
	args := []string{"create"}
	for _, label := range labels {
		args = append(args, "--label", label)
	}
	if topic != "" {
		args = append(args, "--topic", topic)
	}
	args = append(args, title)
	return args
}

func maxTaskID(t *testing.T, req *Request) int {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(req.TskHome, "index"))
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatalf("read index dir: %v", err)
	}
	maxID := 0
	for _, ent := range entries {
		var id int
		if _, err := fmt.Sscanf(ent.Name(), "%d", &id); err == nil && id > maxID {
			maxID = id
		}
	}
	return maxID
}

func createTask(t *testing.T, req *Request, title string, topic string, labels []string) int {
	t.Helper()
	before := maxTaskID(t, req)
	runTskOK(t, req, createArgs(title, topic, labels)...)
	id := maxTaskID(t, req)
	if id <= before {
		t.Fatalf("create task: expected new index entry after create (before=%d after=%d)", before, id)
	}
	req.TaskID = id
	req.Title = title
	req.Topic = topic
	req.Labels = labels
	return id
}

func advanceTask(t *testing.T, req *Request, id int, note string) {
	t.Helper()
	args := []string{"advance", fmt.Sprintf("%d", id)}
	if note != "" {
		args = append(args, "--note", note)
	}
	runTskOK(t, req, args...)
}

func stageTask(t *testing.T, req *Request, id int, stage, note string) {
	t.Helper()
	args := []string{"stage", fmt.Sprintf("%d", id), stage}
	if note != "" {
		args = append(args, "--note", note)
	}
	runTskOK(t, req, args...)
}

func advanceTo(t *testing.T, req *Request, id int, targetStage string) {
	t.Helper()
	order := []string{"create", "in_process", "clarification", "implementation", "verification", "summary"}
	taskDir := findTaskDirByID(t, req, id)
	task := readTaskJSON(t, taskDir)
	cur := task.Stage
	for cur != targetStage {
		switch cur {
		case "create":
			advanceTask(t, req, id, "")
			cur = "in_process"
		case "in_process":
			advanceTask(t, req, id, "")
			cur = "clarification"
		case "clarification":
			runTskOK(t, req, "clarify", "add", fmt.Sprintf("%d", id), "auto question")
			runTskOK(t, req, "clarify", "confirm", fmt.Sprintf("%d", id), "-y")
			cur = "implementation"
		case "implementation":
			advanceTask(t, req, id, "")
			cur = "verification"
		case "verification":
			advanceTask(t, req, id, "")
			cur = "summary"
		default:
			t.Fatalf("advanceTo: cannot reach %q from %q", targetStage, cur)
		}
	}
	_ = order
}

func findTaskDirByID(t *testing.T, req *Request, id int) string {
	t.Helper()
	rel := readIndex(t, req, id)
	abs := taskAbs(req, rel)
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("task dir %s: %v", abs, err)
	}
	return abs
}

func assertErrIsNil(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("%s should exist", path)
	}
}

func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("%s should not exist", path)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Fatalf("%s should exist", path)
	}
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("%s should be a directory", path)
	}
}

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected %q in %q", substr, s)
	}
}

func assertHelpOK(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode != 0 {
		t.Fatalf("help exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("help stderr should be empty, got %q", resp.Stderr)
	}
	if resp.Stdout == "" {
		t.Fatal("help stdout should not be empty")
	}
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("help stdout should end with newline, got %q", resp.Stdout)
	}
}

func assertStderrContainsCount(t *testing.T, stderr, substr string, want int) {
	t.Helper()
	got := strings.Count(stderr, substr)
	if got != want {
		t.Fatalf("stderr substring %q: got %d occurrences want %d; stderr=%q", substr, got, want, stderr)
	}
}

func topLevelSubcommands() []string {
	return []string{
		"create", "list", "show", "status", "advance", "stage", "next",
		"label", "topic", "clarify", "followup", "done",
	}
}

func assertTopHelpStdout(t *testing.T, resp *Response) {
	t.Helper()
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "Usage:")
	for _, cmd := range topLevelSubcommands() {
		assertContains(t, resp.Stdout, cmd)
	}
	assertContains(t, resp.Stdout, "Run tsk <command> --help for command-specific options.")
}

func assertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Fatalf("expected %q not in %q", substr, s)
	}
}

func assertIndexEquals(t *testing.T, req *Request, id int, wantRel string) {
	t.Helper()
	got := readIndex(t, req, id)
	want := filepath.ToSlash(wantRel)
	if got != want {
		t.Fatalf("index/%d: got %q want %q", id, got, want)
	}
}

func assertTaskStage(t *testing.T, req *Request, id int, wantStage string) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	if task.Stage != wantStage {
		t.Fatalf("task %d stage: got %q want %q", id, task.Stage, wantStage)
	}
}

func assertTopicPathNull(t *testing.T, req *Request, id int) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	s := strings.TrimSpace(string(task.TopicPath))
	if s != "null" {
		t.Fatalf("task %d topic_path: got %s want null", id, s)
	}
}

func assertTopicPathEquals(t *testing.T, req *Request, id int, want []string) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	var got []string
	if err := json.Unmarshal(task.TopicPath, &got); err != nil {
		t.Fatalf("task %d topic_path: %v", id, err)
	}
	if len(got) != len(want) {
		t.Fatalf("task %d topic_path: got %v want %v", id, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("task %d topic_path: got %v want %v", id, got, want)
		}
	}
}

func assertLabelsSorted(t *testing.T, req *Request, id int, want []string) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	if len(task.Labels) != len(want) {
		t.Fatalf("task %d labels: got %v want %v", id, task.Labels, want)
	}
	for i := range want {
		if task.Labels[i] != want[i] {
			t.Fatalf("task %d labels: got %v want %v", id, task.Labels, want)
		}
	}
}

func assertStdoutTrimmedEquals(t *testing.T, stdout, want string) {
	t.Helper()
	assert.Output(t, stdout, v2StdoutTemplate(want))
}

func v2StdoutTemplate(body string) string {
	if body == "" {
		return "---\nversion: 2\n---\n"
	}
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return "---\nversion: 2\n---\n" + body
}

func assertEventsCountAtLeast(t *testing.T, req *Request, min int) {
	t.Helper()
	path := filepath.Join(req.TskHome, "events.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read events.jsonl: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if string(data) == "" {
		lines = nil
	}
	if len(lines) < min {
		t.Fatalf("events.jsonl: got %d lines want at least %d", len(lines), min)
	}
}

func assertLastEventCommand(t *testing.T, req *Request, wantCommand string) {
	t.Helper()
	path := filepath.Join(req.TskHome, "events.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read events.jsonl: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		t.Fatal("events.jsonl empty")
	}
	var ev eventLine
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &ev); err != nil {
		t.Fatalf("parse last event: %v", err)
	}
	if ev.Command != wantCommand {
		t.Fatalf("last event command: got %q want %q", ev.Command, wantCommand)
	}
}

func globFollowupFiles(t *testing.T, taskDir string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(taskDir, "context", "followup-*.md"))
	if err != nil {
		t.Fatalf("glob followup files: %v", err)
	}
	return matches
}

func ensureHelpersUsed() {
	_ = findModuleRoot
	_ = getTskBin
	_ = tskEnvBaseDrop
	_ = envKey
	_ = filterEnvKeys
	_ = tskEnv
	_ = captureCommandOutput
	_ = runTskCmd
	_ = runTskOK
	_ = slugify
	_ = taskDirName
	_ = inboxTaskRel
	_ = topicTaskRel
	_ = taskAbs
	_ = readIndex
	_ = readTaskJSON
	_ = readClarifyBatch
	_ = parseCreatedAt
	_ = createArgs
	_ = maxTaskID
	_ = createTask
	_ = advanceTask
	_ = stageTask
	_ = advanceTo
	_ = findTaskDirByID
	_ = assertErrIsNil
	_ = assertFileExists
	_ = assertFileNotExists
	_ = assertDirExists
	_ = assertContains
	_ = assertHelpOK
	_ = assertStderrContainsCount
	_ = topLevelSubcommands
	_ = assertTopHelpStdout
	_ = assertNotContains
	_ = assertIndexEquals
	_ = assertTaskStage
	_ = assertTopicPathNull
	_ = assertTopicPathEquals
	_ = assertLabelsSorted
	_ = assertStdoutTrimmedEquals
	_ = v2StdoutTemplate
	_ = assertEventsCountAtLeast
	_ = assertLastEventCommand
	_ = globFollowupFiles
}
```