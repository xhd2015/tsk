## Expected

- Exit code 0; stdout `2` (id not reused).
- Task 2 exists; index/1 still absent.

## Exit Code

- 0

```go
import "path/filepath"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "2")

	wantRel := inboxTaskRel(2, "create", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, 2, wantRel)
	assertFileNotExists(t, filepath.Join(req.TskHome, "index", "1"))
}
```
