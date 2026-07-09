## Expected

- Exit code 1.
- Stderr mentions invalid transition or stage.
- Task directory remains `inbox/1-create-add-dark-mode/`.
- `index/1` unchanged.
- `task.json` stage remains `create`.

## Errors

- Non-zero exit; error message on stderr.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, stderr=%q", resp.Stderr)
	}
	if resp.Stderr == "" {
		t.Fatal("stderr should contain error message")
	}

	wantRel := inboxTaskRel(req.TaskID, "create", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTaskStage(t, req, req.TaskID, "create")
}
```