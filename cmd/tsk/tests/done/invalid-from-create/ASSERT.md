## Expected

- Exit code 1.
- The task stays at `create` because `--force` was not supplied.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 || !strings.Contains(resp.Stderr, "done only from summary or user_followup") {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertTaskStage(t, req, req.TaskID, "create")
}
```
