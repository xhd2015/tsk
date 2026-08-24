## Expected

- Exit code 1; stderr mentions nested tasks and `--recursive`.
- Parent and child dirs + indexes unchanged.

## Exit Code

- 1

```go
import (
	"path/filepath"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 {
		t.Fatalf("exit=%d want 1 stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "nested tasks") || !strings.Contains(resp.Stderr, "--recursive") {
		t.Fatalf("stderr=%q missing nested/--recursive", resp.Stderr)
	}

	parentRel := inboxTaskRel(req.TaskID, "create", req.Title)
	assertDirExists(t, taskAbs(req, parentRel))
	assertIndexEquals(t, req, req.TaskID, parentRel)

	childID := 2
	childRel := filepath.ToSlash(filepath.Join(parentRel, taskDirName(childID, "create", "child work")))
	assertDirExists(t, taskAbs(req, childRel))
	assertIndexEquals(t, req, childID, childRel)
}
```
