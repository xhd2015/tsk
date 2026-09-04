## Expected

- Exit 0; would-delete parent then child; dirs and indexes remain.

## Exit Code

- 0

```go
import (
	"fmt"
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	want := fmt.Sprintf("[dry-run] would delete %d  [%d] parent work\n[dry-run] would delete 2  [2] child work\n", req.TaskID, req.TaskID)
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}

	parentRel := inboxTaskRel(req.TaskID, req.Title)
	assertDirExists(t, taskAbs(req, parentRel))
	assertIndexEquals(t, req, req.TaskID, parentRel)

	childID := 2
	childRel := filepath.ToSlash(filepath.Join(parentRel, taskDirName(childID, "child work")))
	assertDirExists(t, taskAbs(req, childRel))
	assertIndexEquals(t, req, childID, childRel)
}
```
