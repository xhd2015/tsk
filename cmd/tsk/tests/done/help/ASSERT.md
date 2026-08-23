## Expected

- Exit code 0; stderr empty.
- Help documents `--force`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Usage: tsk done [--force] <id>") || !strings.Contains(resp.Stdout, "--force") {
		t.Fatalf("stdout=%q missing forced completion help", resp.Stdout)
	}
}
```
