## Expected

- Non-zero exit; stderr mentions task(s); topic dir remains.

## Exit Code

- 1

```go
import "os"
import "path/filepath"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected failure")
	}
	assertContains(t, resp.Stderr, "still has")
	dir := filepath.Join(req.TskHome, "topics", "eng")
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("topic dir should remain: %v", err)
	}
}
```
