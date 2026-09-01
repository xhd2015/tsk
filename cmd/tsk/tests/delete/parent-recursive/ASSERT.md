## Expected

- Exit code 0; stdout `deleted <parent>` only.
- Parent and child dirs + indexes gone.

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
	assertStdoutTrimmedEquals(t, resp.Stdout, fmt.Sprintf("deleted %d", req.TaskID))

	parentRel := inboxTaskRel(req.TaskID, req.Title)
	assertFileNotExists(t, taskAbs(req, parentRel))
	assertFileNotExists(t, filepath.Join(req.TskHome, "index", fmt.Sprintf("%d", req.TaskID)))

	childID := 2
	childRel := filepath.ToSlash(filepath.Join(parentRel, taskDirName(childID, "child work")))
	assertFileNotExists(t, taskAbs(req, childRel))
	assertFileNotExists(t, filepath.Join(req.TskHome, "index", fmt.Sprintf("%d", childID)))
}
```
