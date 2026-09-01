## Expected

- Exit code 0; stderr empty.
- Stdout starts with `Usage: tsk tree`.
- Mentions `--json` and inbox root level.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	if !strings.HasPrefix(resp.Stdout, "Usage: tsk tree") {
		t.Fatalf("stdout=%q missing usage prefix", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "--json") {
		t.Fatalf("stdout=%q missing --json", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "project") {
		t.Fatalf("stdout=%q missing project grouping mention", resp.Stdout)
	}
}
```

