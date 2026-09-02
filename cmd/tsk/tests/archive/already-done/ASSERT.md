## Expected

- Exit code 1.
- Stage stays `done`.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 || !strings.Contains(resp.Stderr, "task is already done") {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertTaskStage(t, req, req.TaskID, "done")
}
```
