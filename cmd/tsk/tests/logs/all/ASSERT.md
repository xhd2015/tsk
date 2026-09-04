## Expected

- Exit 0.
- Stdout has `add` then `list`, footer `2 logs`. No `logs` self-row.

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
	if resp.Stderr != "" {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
	addLine := fmt.Sprintf("2026-07-09T12:00:00+08:00  ok  add  task=%d", req.TaskID)
	listLine := "2026-07-09T12:00:00+08:00  ok  list"
	want := addLine + "\n" + listLine + "\n2 logs\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
	if strings.Contains(resp.Stdout, "  logs") {
		t.Fatalf("logs should not record itself: %q", resp.Stdout)
	}
}
```
