## Expected

- Exit code 0; stderr empty.
- The done task leaf is present without ANSI sequences.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, "[1] finished  (done)") {
		t.Fatalf("stdout=%q missing done task leaf", resp.Stdout)
	}
}
```
