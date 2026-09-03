## Expected

- Exit 0; frontmatter `name: tsk/working-on-task`.
- Body covers expanded intake, kck poll commands, and independent verify.

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
	if !strings.Contains(out, "name: tsk/working-on-task") {
		t.Fatalf("expected name: tsk/working-on-task:\n%s", out)
	}
	for _, want := range []string{
		"Intent",
		"Project",
		"Prerequisites",
		"Definition of done",
		"Verification steps",
		"Artifacts needed",
		"kck grok new",
		"kck grok messages",
		"kck grok snapshot",
		"Independent verify",
		"Report",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
}
```
