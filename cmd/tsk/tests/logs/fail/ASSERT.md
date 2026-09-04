## Expected

- Exit 0.
- Three mutations: `add` ok, `done` ok, `done` fail; footer `3 logs`.

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
	id := req.TaskID
	want := strings.Join([]string{
		fmt.Sprintf("2026-07-09T12:00:00+08:00  ok  add  task=%d", id),
		fmt.Sprintf("2026-07-09T12:00:00+08:00  ok  done  task=%d  stage=done", id),
		fmt.Sprintf("2026-07-09T12:00:00+08:00  fail  done  task=%d", id),
		"3 logs",
		"",
	}, "\n")
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
