## Expected

- Parent moves to `inbox/[1]-in_process-add-dark-mode/`.
- Child remains nested; `index/2` uses the new parent dirname prefix.

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

	oldParent := inboxTaskRel(1, "create", "add dark mode")
	assertFileNotExists(t, taskAbs(req, oldParent))

	newParent := inboxTaskRel(1, "in_process", "add dark mode")
	assertDirExists(t, taskAbs(req, newParent))
	assertIndexEquals(t, req, 1, newParent)

	childRel := filepath.ToSlash(filepath.Join(newParent, taskDirName(2, "create", "child detail")))
	assertDirExists(t, taskAbs(req, childRel))
	assertIndexEquals(t, req, 2, childRel)
}
```
