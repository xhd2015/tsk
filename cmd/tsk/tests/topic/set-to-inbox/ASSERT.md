## Expected

- Exit code 0.
- Task directory under `inbox/1-create-return-home/`.
- `topic_path` is null in `task.json`.
- `index/1` points at inbox path.

## Side Effects

- Topic task dir removed.

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

	oldRel := topicTaskRel(req.Topic, req.TaskID, "create", req.Title)
	assertFileNotExists(t, taskAbs(req, oldRel))

	wantRel := inboxTaskRel(req.TaskID, "create", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTopicPathNull(t, req, req.TaskID)
}
```