## Expected

- Exit code 0.
- Task directory moved from inbox to `topics/eng/backend/1-create-move-me/`.
- `index/1` updated to topic-relative path.
- `task.json` `topic_path` matches `["eng","backend"]`.

## Side Effects

- Inbox task dir removed; topic dir contains task.

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

	oldRel := inboxTaskRel(req.TaskID, "create", req.Title)
	assertFileNotExists(t, taskAbs(req, oldRel))

	wantRel := topicTaskRel(req.Topic, req.TaskID, "create", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTopicPathEquals(t, req, req.TaskID, strings.Split(req.Topic, "/"))
}
```