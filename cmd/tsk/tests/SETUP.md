# Scenario

**Feature**: isolated TSK_HOME per test; build tsk binary once per doctest session

```
# temp WorkRoot + {WorkRoot}/.tsk per leaf; session-cached tsk binary
tsk <subcommand> -> stdout/stderr + filesystem side effects under TSK_HOME
```

## Preconditions

- The tsk Go module root is two levels above the test tree (`github.com/xhd2015/tsk`).
- Go toolchain is available on PATH.
- Session cache: `{DOCTEST_FIXTURE_ROOT or ~/Library/Caches/doctest/fixtures}/{DOCTEST_SESSION_ID}/bin/tsk` (file-locked build).

## Context

Each leaf runs `tsk` with `TSK_HOME={WorkRoot}/.tsk` and `TSK_DATE=2026-07-09`. Helper `runTskCmd` runs additional CLI invocations during `Setup` chains. Per-leaf state is never shared across leaves.

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
	"unicode"

	"github.com/xhd2015/doctest/assert"
)

const tskDate = "2026-07-09"

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

func fixtureCacheBase(t *testing.T) string {
	t.Helper()
	base := os.Getenv("DOCTEST_FIXTURE_ROOT")
	if base != "" {
		return base
	}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(home, "Library", "Caches", "doctest", "fixtures")
}

func fixtureSessionRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(fixtureCacheBase(t), DOCTEST_SESSION_ID)
}

func sessionTskBin(t *testing.T) string {
	t.Helper()
	return filepath.Join(fixtureSessionRoot(t), "bin", "tsk")
}

func withFlock(t *testing.T, lockPath string, fn func()) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		t.Fatalf("mkdir lock dir: %v", err)
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("open lock %s: %v", lockPath, err)
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		t.Fatalf("flock %s: %v", lockPath, err)
	}
	defer func() { _ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN) }()
	fn()
}

func getTskBin(t *testing.T) string {
	t.Helper()
	bin := sessionTskBin(t)
	if _, err := os.Stat(bin); err == nil {
		return bin
	}
	lockPath := filepath.Join(fixtureSessionRoot(t), "bin", ".lock")
	withFlock(t, lockPath, func() {
		if _, err := os.Stat(bin); err == nil {
			return
		}
		modRoot := filepath.Dir(filepath.Dir(DOCTEST_ROOT))
		if modRoot == "" {
			modRoot = findModuleRoot(DOCTEST_ROOT)
		}
		cmd := exec.Command("go", "build", "-o", bin, "./cmd/tsk")
		cmd.Dir = modRoot
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("build tsk: %v\n%s", err, out)
		}
	})
	return bin
}

func Setup(t *testing.T, req *Request) error {
	workRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		return fmt.Errorf("resolve work root: %w", err)
	}
	req.WorkRoot = workRoot
	req.TskHome = filepath.Join(req.WorkRoot, ".tsk")
	ensureHelpersUsed()
	return os.MkdirAll(req.TskHome, 0o755)
}

func tskEnv(req *Request) []string {
	return append(os.Environ(),
		"TSK_HOME="+req.TskHome,
		"TSK_DATE="+tskDate,
	)
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
	_ = fixtureCacheBase
	_ = fixtureSessionRoot
	_ = sessionTskBin
	_ = withFlock
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