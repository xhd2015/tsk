## Expected

- Exit 0; stdout `[dry-run] would delete 1  [1] oops`.
- Task dir and index still present; `show` succeeds.

## Exit Code

- 0

```go
import "fmt"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	want := fmt.Sprintf("[dry-run] would delete %d  [%d] %s\n", req.TaskID, req.TaskID, req.Title)
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}

	wantRel := inboxTaskRel(req.TaskID, req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)

	show := runTskCmd(t, req, "show", fmt.Sprintf("%d", req.TaskID))
	if show.ExitCode != 0 {
		t.Fatalf("show after dry-run: exit=%d stderr=%q", show.ExitCode, show.Stderr)
	}
}
```
