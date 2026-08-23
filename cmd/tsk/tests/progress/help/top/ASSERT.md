## Expected

- Exit code 0; stderr empty.
- Stdout starts with `Usage: tsk progress`.
- Mentions `add`, `list`, `edit`, `archive`, `show`, and `--status`.

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
	if !strings.HasPrefix(resp.Stdout, "Usage: tsk progress") {
		t.Fatalf("stdout=%q missing usage prefix", resp.Stdout)
	}
	for _, sub := range []string{"add", "list", "edit", "archive", "show", "--status"} {
		if !strings.Contains(resp.Stdout, sub) {
			t.Fatalf("stdout=%q missing %q", resp.Stdout, sub)
		}
	}
}
```
