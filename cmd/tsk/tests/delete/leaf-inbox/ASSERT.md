## Expected

- Exit code 0; stdout `deleted <id>`; stderr empty.
- Task directory and `index/<id>` gone.
- `show <id>` fails with not found.

## Exit Code

- 0

```go
import (
	"fmt"
	"path/filepath"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, fmt.Sprintf("deleted %d", req.TaskID))

	wantRel := inboxTaskRel(req.TaskID, req.Title)
	assertFileNotExists(t, taskAbs(req, wantRel))
	assertFileNotExists(t, filepath.Join(req.TskHome, "index", fmt.Sprintf("%d", req.TaskID)))

	show := runTskCmd(t, req, "show", fmt.Sprintf("%d", req.TaskID))
	if show.ExitCode == 0 || !strings.Contains(show.Stderr, "not found") {
		t.Fatalf("show after delete: exit=%d stderr=%q", show.ExitCode, show.Stderr)
	}
}
```
