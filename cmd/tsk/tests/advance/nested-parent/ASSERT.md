## Expected

- Parent directory basename stays `inbox/[1]-add-dark-mode/`.
- Child remains nested under that parent; indexes unchanged.
- Parent stage becomes `in_process`.

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

	parentRel := inboxTaskRel(1, "add dark mode")
	assertDirExists(t, taskAbs(req, parentRel))
	assertIndexEquals(t, req, 1, parentRel)
	assertTaskStage(t, req, 1, "in_process")

	childRel := filepath.ToSlash(filepath.Join(parentRel, taskDirName(2, "child detail")))
	assertDirExists(t, taskAbs(req, childRel))
	assertIndexEquals(t, req, 2, childRel)
}
```
