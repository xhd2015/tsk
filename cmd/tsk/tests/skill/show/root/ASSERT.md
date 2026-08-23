## Expected

- Exit 0.
- Frontmatter `name: tsk`; retrieve example `tsk skill --show`.
- No `--cursor` / `--global` install plumbing in body.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	out := resp.Stdout
	if !strings.Contains(out, "name: tsk") {
		t.Fatalf("expected name: tsk in stdout:\n%s", out)
	}
	if !strings.Contains(out, "tsk skill --show") {
		t.Fatalf("expected retrieve example in stdout:\n%s", out)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "topic") {
		t.Fatalf("expected multi-topic index language:\n%s", out)
	}
	if strings.Contains(out, "--cursor") || strings.Contains(out, "--global") {
		t.Fatalf("root skill must not document install plumbing flags:\n%s", out)
	}
}
```
