## Expected

- Default: only open task; hides done/archived.
- `--done`: only done.
- `--archived`: only archived.
- `--done --archived`: both terminals, no open.
- `--all`: all three stages.
- `--stage` with `--done` errors.

## Exit Code

- 0

```go
import (
	"fmt"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "dot-pkgs  github.com/xhd2015/dot-pkgs")
	assertContains(t, resp.Stdout, fmt.Sprintf("[%d] active-me  (create)", req.OpenID))
	assertNotContains(t, resp.Stdout, "finish-me")
	assertNotContains(t, resp.Stdout, "shelve-me")
	assertContains(t, resp.Stdout, "1 task, 1 project")

	onlyDone := runTskOK(t, req, "project", "tree", "--done", "--plain")
	assertContains(t, onlyDone.Stdout, fmt.Sprintf("[%d] finish-me  (done)", req.TaskID))
	assertNotContains(t, onlyDone.Stdout, "active-me")
	assertNotContains(t, onlyDone.Stdout, "shelve-me")
	assertContains(t, onlyDone.Stdout, "1 task, 1 project")

	onlyArch := runTskOK(t, req, "project", "tree", "--archived", "--plain")
	assertContains(t, onlyArch.Stdout, fmt.Sprintf("[%d] shelve-me  (archived)", req.ArchivedID))
	assertNotContains(t, onlyArch.Stdout, "active-me")
	assertNotContains(t, onlyArch.Stdout, "finish-me")
	assertContains(t, onlyArch.Stdout, "1 task, 1 project")

	bothTerm := runTskOK(t, req, "project", "tree", "--done", "--archived", "--plain")
	assertContains(t, bothTerm.Stdout, fmt.Sprintf("[%d] finish-me  (done)", req.TaskID))
	assertContains(t, bothTerm.Stdout, fmt.Sprintf("[%d] shelve-me  (archived)", req.ArchivedID))
	assertNotContains(t, bothTerm.Stdout, "active-me")
	assertContains(t, bothTerm.Stdout, "2 tasks, 1 project")

	allStages := runTskOK(t, req, "project", "tree", "--all", "--plain")
	assertContains(t, allStages.Stdout, "active-me")
	assertContains(t, allStages.Stdout, "finish-me")
	assertContains(t, allStages.Stdout, "shelve-me")
	assertContains(t, allStages.Stdout, "3 tasks, 1 project")

	conflict := runTskCmd(t, req, "project", "tree", "--stage", "done", "--done")
	if conflict.ExitCode != 1 || !strings.Contains(conflict.Stderr, "--stage conflicts with --done/--archived") {
		t.Fatalf("conflict exit=%d stderr=%q", conflict.ExitCode, conflict.Stderr)
	}
}
```
