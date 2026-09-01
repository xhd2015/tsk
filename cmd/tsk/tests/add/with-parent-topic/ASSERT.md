## Expected

- Child nested under topic parent; inherits `topic_path: ["kb"]`.

## Exit Code

- 0

```go
import (
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "2")

	parentRel := topicTaskRel("kb", 1, "parent report")
	childRel := filepath.ToSlash(filepath.Join(parentRel, taskDirName(2, "child detail")))
	assertDirExists(t, taskAbs(req, childRel))
	assertIndexEquals(t, req, 2, childRel)
	task := readTaskJSON(t, taskAbs(req, childRel))
	if task.ParentID != 1 {
		t.Fatalf("parent_id: got %d want 1", task.ParentID)
	}
	assertTopicPathEquals(t, req, 2, []string{"kb"})
}
```
