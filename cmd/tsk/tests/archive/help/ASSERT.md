## Expected

- Exit code 0; stderr empty.
- Help documents any non-terminal stage and `--force` compatibility.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Usage: tsk archive [--force] <id>") ||
		!strings.Contains(resp.Stdout, "any non-terminal stage") ||
		!strings.Contains(resp.Stdout, "--force") {
		t.Fatalf("stdout=%q missing archive help", resp.Stdout)
	}
}
```
