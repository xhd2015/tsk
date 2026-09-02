## Expected

- Exit code 1.
- Stage stays `archived`.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 || !strings.Contains(resp.Stderr, "task is already archived") {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertTaskStage(t, req, req.TaskID, "archived")
}
```
