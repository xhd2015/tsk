## Expected

- Exit 0; task in inbox; topic_path null.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	oldRel := topicTaskRel(req.Topic, req.TaskID, req.Title)
	assertFileNotExists(t, taskAbs(req, oldRel))
	wantRel := inboxTaskRel(req.TaskID, req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTopicPathNull(t, req, req.TaskID)
}
```
