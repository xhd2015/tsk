## Expected

- Exit code 1.
- Stderr contains `Error:` and `task not found`.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 {
		t.Fatalf("exit %d, want 1", resp.ExitCode)
	}
	if !strings.Contains(resp.Stderr, "Error:") {
		t.Fatalf("stderr=%q missing Error:", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "task not found") {
		t.Fatalf("stderr=%q missing task not found", resp.Stderr)
	}
}
```
