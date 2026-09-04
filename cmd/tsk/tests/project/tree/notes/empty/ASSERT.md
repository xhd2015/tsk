## Expected

- Exit 0; task leaf only; no `notes` group.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "[1] one")
	if strings.Contains(resp.Stdout, "\n    ├── notes\n") || strings.Contains(resp.Stdout, "\n        notes\n") {
		t.Fatalf("unexpected notes node: %q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "1 task, 1 project")
}
```
