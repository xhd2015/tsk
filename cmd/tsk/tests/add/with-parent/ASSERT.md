## Expected

- Exit code 0; stdout child id `2`.
- Child dir nested under parent: `inbox/[1]-parent-work/[2]-child-work/`.
- `index/2` points at nested path; `parent_id: 2`'s task.json has `parent_id: 1`, `topic_path: null`.

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

	parentRel := inboxTaskRel(1, "parent work")
	childRel := filepath.ToSlash(filepath.Join(parentRel, taskDirName(2, "child work")))
	assertDirExists(t, taskAbs(req, childRel))
	assertIndexEquals(t, req, 2, childRel)

	task := readTaskJSON(t, taskAbs(req, childRel))
	if task.ParentID != 1 {
		t.Fatalf("parent_id: got %d want 1", task.ParentID)
	}
	assertTopicPathNull(t, req, 2)
}
```
