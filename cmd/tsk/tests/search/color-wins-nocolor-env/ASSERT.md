## Expected

- Exit code 0.
- ANSI present despite `NO_COLOR=1` because `--color` forces Always.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "\x1b[") {
		t.Fatalf("expected ANSI with --color despite NO_COLOR, stdout=%q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "env-color-token")
}
```
