## Expected

- Exit 1; stderr contains `not found`.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 || !strings.Contains(resp.Stderr, "not found") {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
}
```