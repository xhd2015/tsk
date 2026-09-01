## Expected

- Exit 0; stdout `removed topic eng`; topics/eng gone.

## Exit Code

- 0

```go
import "os"
import "path/filepath"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "removed topic eng\n")
	dir := filepath.Join(req.TskHome, "topics", "eng")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("topic dir should be gone: %v", err)
	}
}
```
