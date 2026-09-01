## Expected

- Default list: project branch present, `0 tasks`, no done leaf.
- `--stage done --plain`: shows the done task leaf.

## Exit Code

- 0

```go
import "fmt"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "dot-pkgs  github.com/xhd2015/dot-pkgs")
	assertContains(t, resp.Stdout, "0 tasks, 1 project")
	assertNotContains(t, resp.Stdout, "finish-me")

	doneList := runTskOK(t, req, "project", "tree", "--stage", "done", "--plain")
	assertContains(t, doneList.Stdout, fmt.Sprintf("[%d] finish-me  (done)", req.TaskID))
	assertContains(t, doneList.Stdout, "1 task, 1 project")
}
```
