# Scenario

**Feature**: isolated TSK_HOME per test; build check-channel-activity binary once per doctest session

```
# temp WorkRoot + {WorkRoot}/.tsk per leaf; session-cached binary
check-channel-activity --channel-id ID --exec-if-idle-1h "LINE" -> shellwords parse -> exec + state
```

## Preconditions

- Module root `github.com/xhd2015/tsk` is three levels above the test tree (`DOCTEST_ROOT/../../..`).
- Go toolchain on PATH.
- Session cache: `{DOCTEST_FIXTURE_ROOT or ~/Library/Caches/doctest/fixtures}/{DOCTEST_SESSION_ID}/bin/check-channel-activity` (file-locked build).
- Child env sets `TSK_HOME={WorkRoot}/.tsk`; strips parent `TSK_HOME` to avoid leakage.

## Context

Each leaf runs `check-channel-activity` with an isolated `TSK_HOME`. Channel fixtures
are written directly under `{WorkRoot}/.tsk/channels/`. Exec verification uses a shell
script that touches `{WorkRoot}/notify.marker`.

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

const (
	channelID      = "eng-alerts"
	channelName    = "Eng Alerts"
	oldActivityTS  = "2026-07-09T01:00:00Z"
	oldCreatedAtTS = "2026-07-09T00:30:00Z"
)

type channelParticipant struct {
	Handle   string `json:"handle"`
	JoinedAt string `json:"joined_at"`
}

type channelJSON struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Status       string               `json:"status"`
	Participants []channelParticipant `json:"participants"`
	CreatedAt    string               `json:"created_at"`
	UpdatedAt    string               `json:"updated_at"`
}

type channelMessage struct {
	ID        int    `json:"id"`
	Sender    string `json:"sender"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type notifyState struct {
	ChannelID      string `json:"channel_id"`
	LastActivityAt string `json:"last_activity_at"`
	LastNotifiedAt string `json:"last_notified_at"`
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

func sessionCheckBin(t *testing.T) string {
	t.Helper()
	return filepath.Join(fixtureSessionRoot(t), "bin", "check-channel-activity")
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

func getCheckBin(t *testing.T) string {
	t.Helper()
	bin := sessionCheckBin(t)
	if _, err := os.Stat(bin); err == nil {
		return bin
	}
	lockPath := filepath.Join(fixtureSessionRoot(t), "bin", ".check-channel-activity.lock")
	withFlock(t, lockPath, func() {
		if _, err := os.Stat(bin); err == nil {
			return
		}
		modRoot := findModuleRoot(DOCTEST_ROOT)
		if modRoot == "" {
			t.Fatal("find module root")
		}
		cmd := exec.Command("go", "build", "-o", bin, "./script/check-channel-activity")
		cmd.Dir = modRoot
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("build check-channel-activity: %v\n%s", err, out)
		}
		warm := exec.Command(bin, "--help")
		warm.Env = append(os.Environ(), "CHECK_CHANNEL_ACTIVITY_WARMED=1")
		warm.Stdout = nil
		warm.Stderr = nil
		_ = warm.Run()
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
	req.ChannelID = channelID
	req.MarkerPath = filepath.Join(req.WorkRoot, "notify.marker")
	req.ExecScript = filepath.Join(req.WorkRoot, "touch-marker.sh")
	if err := os.MkdirAll(req.TskHome, 0o755); err != nil {
		return err
	}
	if err := writeExecScript(req); err != nil {
		return err
	}
	ensureCheckHelpersUsed()
	return nil
}

var checkEnvDrop = map[string]struct{}{
	"TSK_HOME": {},
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

func checkEnv(req *Request) []string {
	drop := make(map[string]struct{}, len(checkEnvDrop)+len(req.ExtraEnv))
	for k := range checkEnvDrop {
		drop[k] = struct{}{}
	}
	for _, e := range req.ExtraEnv {
		drop[envKey(e)] = struct{}{}
	}
	env := filterEnvKeys(os.Environ(), drop)
	env = append(env, "TSK_HOME="+req.TskHome)
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

func channelsRoot(req *Request) string {
	return filepath.Join(req.TskHome, "channels")
}

func statePath(req *Request) string {
	return filepath.Join(req.TskHome, "channels", "state", req.ChannelID+".json")
}

func writeExecScript(req *Request) error {
	body := "#!/bin/sh\nset -eu\ntouch \"" + req.MarkerPath + "\"\n"
	if err := os.WriteFile(req.ExecScript, []byte(body), 0o755); err != nil {
		return err
	}
	return nil
}

func writeChannelIndex(req *Request, rel string) error {
	dir := filepath.Join(channelsRoot(req), "index")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, req.ChannelID)
	return os.WriteFile(path, []byte(rel+"\n"), 0o644)
}

func writeChannelJSON(channelDir string, ch channelJSON) error {
	if err := os.MkdirAll(channelDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ch, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := filepath.Join(channelDir, "channel.json.tmp")
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(channelDir, "channel.json"))
}

func writeMessages(channelDir string, msgs []channelMessage) error {
	if err := os.MkdirAll(channelDir, 0o755); err != nil {
		return err
	}
	var b strings.Builder
	for _, m := range msgs {
		line, err := json.Marshal(m)
		if err != nil {
			return err
		}
		b.Write(line)
		b.WriteByte('\n')
	}
	path := filepath.Join(channelDir, "messages.jsonl")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return err
	}
	counter := 0
	if len(msgs) > 0 {
		counter = msgs[len(msgs)-1].ID
	}
	return os.WriteFile(filepath.Join(channelDir, "msg-counter"), []byte(fmt.Sprintf("%d\n", counter)), 0o644)
}

func writeActiveChannel(t *testing.T, req *Request, createdAt string, msgs []channelMessage) string {
	t.Helper()
	rel := "active/" + req.ChannelID
	dir := filepath.Join(channelsRoot(req), "active", req.ChannelID)
	ch := channelJSON{
		ID:     req.ChannelID,
		Name:   channelName,
		Status: "active",
		Participants: []channelParticipant{
			{Handle: "agent", JoinedAt: createdAt},
			{Handle: "alice", JoinedAt: createdAt},
		},
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	if err := writeChannelIndex(req, rel); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := writeChannelJSON(dir, ch); err != nil {
		t.Fatalf("write channel.json: %v", err)
	}
	if err := writeMessages(dir, msgs); err != nil {
		t.Fatalf("write messages: %v", err)
	}
	last := createdAt
	if len(msgs) > 0 {
		last = msgs[len(msgs)-1].CreatedAt
	}
	return last
}

func writeArchivedChannel(t *testing.T, req *Request) {
	t.Helper()
	rel := "archive/" + req.ChannelID
	dir := filepath.Join(channelsRoot(req), "archive", req.ChannelID)
	ch := channelJSON{
		ID:     req.ChannelID,
		Name:   channelName,
		Status: "archived",
		Participants: []channelParticipant{
			{Handle: "agent", JoinedAt: oldCreatedAtTS},
			{Handle: "alice", JoinedAt: oldCreatedAtTS},
		},
		CreatedAt: oldCreatedAtTS,
		UpdatedAt: oldActivityTS,
	}
	if err := writeChannelIndex(req, rel); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := writeChannelJSON(dir, ch); err != nil {
		t.Fatalf("write channel.json: %v", err)
	}
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "old", CreatedAt: oldActivityTS,
	}}
	if err := writeMessages(dir, msgs); err != nil {
		t.Fatalf("write messages: %v", err)
	}
}

func writeNotifyState(t *testing.T, req *Request, lastActivity, lastNotified string) {
	t.Helper()
	dir := filepath.Join(req.TskHome, "channels", "state")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir state: %v", err)
	}
	st := notifyState{
		ChannelID:      req.ChannelID,
		LastActivityAt: lastActivity,
		LastNotifiedAt: lastNotified,
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		t.Fatalf("marshal state: %v", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(statePath(req), data, 0o644); err != nil {
		t.Fatalf("write state: %v", err)
	}
}

func recentActivityTS() string {
	return time.Now().UTC().Add(-30 * time.Minute).Format(time.RFC3339)
}

func defaultExecLine(script string) string {
	return fmt.Sprintf("/bin/sh %q", script)
}

func defaultCheckArgs(req *Request, extra ...string) []string {
	base := []string{
		"--channel-id", req.ChannelID,
		"--exec-if-idle-1h", defaultExecLine(req.ExecScript),
	}
	return append(base, extra...)
}

func runCheckCmd(t *testing.T, req *Request, args ...string) *Response {
	t.Helper()
	bin := getCheckBin(t)
	cmd := exec.Command(bin, args...)
	cmd.Dir = req.WorkRoot
	cmd.Env = checkEnv(req)
	resp, err := captureCommandOutput(cmd)
	if err != nil {
		t.Fatalf("run check-channel-activity %v: %v", args, err)
	}
	return resp
}

func runCheckOK(t *testing.T, req *Request, args ...string) *Response {
	t.Helper()
	resp := runCheckCmd(t, req, args...)
	if resp.ExitCode != 0 {
		t.Fatalf("check-channel-activity %v: exit %d stderr=%q", args, resp.ExitCode, resp.Stderr)
	}
	return resp
}

func runWithSIGINT(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	bin := getCheckBin(t)
	cmd := exec.Command(bin, req.Args...)
	cmd.Dir = req.WorkRoot
	cmd.Env = checkEnv(req)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	time.Sleep(100 * time.Millisecond)
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		_ = cmd.Process.Kill()
		return nil, err
	}
	err := cmd.Wait()
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

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected %q in %q", substr, s)
	}
}

func assertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Fatalf("expected %q not in %q", substr, s)
	}
}

func assertStderrErrorPrefix(t *testing.T, stderr string) {
	t.Helper()
	if !strings.Contains(stderr, "Error:") {
		t.Fatalf("stderr should contain Error:, got %q", stderr)
	}
}

func assertStdoutEndsWithNewline(t *testing.T, stdout string) {
	t.Helper()
	if stdout != "" && !strings.HasSuffix(stdout, "\n") {
		t.Fatalf("stdout should end with newline, got %q", stdout)
	}
}

func assertStatusBlock(t *testing.T, stdout, wantStatus, wantLastActivity string) {
	t.Helper()
	assertContains(t, stdout, "channel: "+channelID)
	assertContains(t, stdout, "last_activity: "+wantLastActivity)
	assertContains(t, stdout, "idle:")
	assertContains(t, stdout, "status: "+wantStatus)
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

func ensureCheckHelpersUsed() {
	_ = findModuleRoot
	_ = fixtureCacheBase
	_ = fixtureSessionRoot
	_ = sessionCheckBin
	_ = withFlock
	_ = getCheckBin
	_ = checkEnvDrop
	_ = envKey
	_ = filterEnvKeys
	_ = checkEnv
	_ = captureCommandOutput
	_ = channelsRoot
	_ = statePath
	_ = writeExecScript
	_ = writeChannelIndex
	_ = writeChannelJSON
	_ = writeMessages
	_ = writeActiveChannel
	_ = writeArchivedChannel
	_ = writeNotifyState
	_ = recentActivityTS
	_ = defaultExecLine
	_ = defaultCheckArgs
	_ = runCheckCmd
	_ = runCheckOK
	_ = runWithSIGINT
	_ = assertErrIsNil
	_ = assertFileExists
	_ = assertFileNotExists
	_ = assertContains
	_ = assertNotContains
	_ = assertStderrErrorPrefix
	_ = assertStdoutEndsWithNewline
	_ = assertStatusBlock
	_ = assertHelpOK
}
```