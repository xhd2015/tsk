## Expected

- Exit 0; frontmatter `name: tsk/working-on-task`.
- Body uses numbered subsections (Pick a task via show or
  `tsk project tree --all` + current-project preference), reuse-or-create
  worktree, kck new/watch/send/wait, tsk stages/must-notes, verify, and report.

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
		"## 1. Pick a task",
		"tsk project tree --all",
		"tsk project which",
		"Project",
		"Intent kind",
		"E2E acceptance",
		"Must-notes",
		"tsk show",
		"linked git worktree",
		"reused existing worktree",
		"wrk <project-location> --no-config",
		"wrk --bring",
		"wrk --done",
		"kck grok new",
		"kck grok messages",
		"kck grok snapshot",
		"kck grok send",
		"kck grok wait",
		"tsk done",
		"Report",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
}
```
