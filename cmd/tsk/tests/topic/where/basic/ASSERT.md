## Expected

- Exit code 0; stderr empty.
- Stdout is the absolute topic dir plus newline.
- Path ends with `/topics/knowledge-base`.

## Exit Code

- 0

```go
import (
	"path/filepath"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	want := filepath.Join(req.TskHome, "topics", "knowledge-base") + "\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
	if !strings.Contains(filepath.ToSlash(resp.Stdout), "/topics/knowledge-base") {
		t.Fatalf("stdout missing topics/knowledge-base: %q", resp.Stdout)
	}
}
```
