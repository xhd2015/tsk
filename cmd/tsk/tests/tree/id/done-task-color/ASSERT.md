## Expected

- Exit code 0; stderr empty.
- The done task leaf is gray + struck through; topic and branch characters are unstyled.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}

	ansi, reset := "\x1b[90m\x1b[9m", "\x1b[0m"
	if !strings.Contains(resp.Stdout, ansi+"[1]-done-finished  task 1  done"+reset) {
		t.Fatalf("stdout=%q missing styled task leaf", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, ansi+"projects") || strings.Contains(resp.Stdout, ansi+"└──") {
		t.Fatalf("topic and branch art must not be styled: %q", resp.Stdout)
	}
}
```
