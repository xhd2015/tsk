## Expected

- Exit 0; stderr empty.
- One `add` mutation line with `task=<id>`; no `list`; footer `1 log`.

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
	wantLine := fmt.Sprintf("2026-07-09T12:00:00+08:00  ok  add  task=%d", req.TaskID)
	if resp.Stdout != wantLine+"\n1 log\n" {
		t.Fatalf("stdout=%q want %q", resp.Stdout, wantLine+"\n1 log\n")
	}
	if strings.Contains(resp.Stdout, "  list") {
		t.Fatalf("default logs should hide list: %q", resp.Stdout)
	}
}
```
