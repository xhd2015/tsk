## Expected

- Exit code 0; stderr empty.
- Help documents `--recursive` and contrasts with `done`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Usage: tsk delete [--recursive] [--dry-run] <id>") || !strings.Contains(resp.Stdout, "--recursive") || !strings.Contains(resp.Stdout, "--dry-run") {
		t.Fatalf("stdout=%q missing delete help", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Unlike tsk done") {
		t.Fatalf("stdout=%q should contrast with done", resp.Stdout)
	}
}
```
