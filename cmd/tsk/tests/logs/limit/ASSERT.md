## Expected

- Exit 0.
- Only the second add (`task=<id2>`); footer `1 log`.

## Exit Code

- 0

```go
import "fmt"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	want := fmt.Sprintf("2026-07-09T12:00:00+08:00  ok  add  task=%d\n1 log\n", req.TaskID)
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
