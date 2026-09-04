## Expected

- Exit 0.
- Lines: `add` then `note.add` with `task=`; no `note.list`; footer `2 logs`.

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
	addLine := fmt.Sprintf("2026-07-09T12:00:00+08:00  ok  add  task=%d", req.TaskID)
	noteLine := fmt.Sprintf("2026-07-09T12:00:00+08:00  ok  note.add  task=%d", req.TaskID)
	want := addLine + "\n" + noteLine + "\n2 logs\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
	if strings.Contains(resp.Stdout, "note.list") {
		t.Fatalf("note.list should be hidden: %q", resp.Stdout)
	}
}
```
