## Expected

- Exit code 0; stderr empty.
- Stage `done`; directory renamed to `*-done-*`.
- `index/<id>` updated and stage history records the direct completion.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}

	wantRel := inboxTaskRel(req.TaskID, "done", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTaskStage(t, req, req.TaskID, "done")
}
```
