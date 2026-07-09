## Expected

- Exit code 0.
- Stage `user_followup`; directory renamed to `*-user_followup-*`.
- Exactly one `context/followup-*.md` file containing the message text.
- `index/<id>` updated.

## Side Effects

- New markdown artifact under `context/`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	wantRel := inboxTaskRel(req.TaskID, "user_followup", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTaskStage(t, req, req.TaskID, "user_followup")

	taskDir := findTaskDirByID(t, req, req.TaskID)
	files := globFollowupFiles(t, taskDir)
	if len(files) != 1 {
		t.Fatalf("expected 1 followup file, got %d: %v", len(files), files)
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("read followup file: %v", err)
	}
	if !strings.Contains(string(data), req.Message) {
		t.Fatalf("followup file missing message %q: %s", req.Message, data)
	}
}
```