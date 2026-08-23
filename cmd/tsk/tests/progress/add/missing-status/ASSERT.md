## Expected

- Exit code 1.
- Stderr contains `Error:` and `--status required`.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 {
		t.Fatalf("exit %d, want 1", resp.ExitCode)
	}
	if !strings.Contains(resp.Stderr, "Error:") || !strings.Contains(resp.Stderr, "--status required") {
		t.Fatalf("stderr=%q missing status error", resp.Stderr)
	}
}
```
