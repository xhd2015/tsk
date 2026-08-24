## Expected

- Exit code 0; stdout `deleted <child>`.
- Child dir + index gone; parent dir + index remain.

## Exit Code

- 0

```go
import (
	"fmt"
	"path/filepath"
	"strconv"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, fmt.Sprintf("deleted %d", req.TaskID))

	parentID, perr := strconv.Atoi(req.Message)
	if perr != nil {
		t.Fatalf("parent id: %v", perr)
	}
	parentRel := inboxTaskRel(parentID, "create", "parent work")
	assertDirExists(t, taskAbs(req, parentRel))
	assertIndexEquals(t, req, parentID, parentRel)

	childRel := filepath.ToSlash(filepath.Join(parentRel, taskDirName(req.TaskID, "create", req.Title)))
	assertFileNotExists(t, taskAbs(req, childRel))
	assertFileNotExists(t, filepath.Join(req.TskHome, "index", fmt.Sprintf("%d", req.TaskID)))
}
```
