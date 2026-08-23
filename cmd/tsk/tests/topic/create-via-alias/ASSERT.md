## Expected

- Exit code 0.
- Task directory is under `topics/knowledge-base/`.
- `task.json` `topic_path` is `["knowledge-base"]`, not `["知识库"]`.
- `index/1` points at the canonical relative path.

## Exit Code

- 0

```go
import (
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	wantRel := topicTaskRel("knowledge-base", 1, "create", "x")
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, 1, wantRel)
	task := readTaskJSON(t, taskAbs(req, wantRel))
	assertTopicPathEquals(t, req, 1, []string{"knowledge-base"})
	if strings.Contains(string(task.TopicPath), "知识库") {
		t.Fatalf("topic_path should be canonical, got %s", task.TopicPath)
	}
}
```
