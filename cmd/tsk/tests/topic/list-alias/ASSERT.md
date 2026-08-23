## Expected

- Exit code 0; stderr empty.
- Stdout is the created task id.

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
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	got := strings.TrimSpace(resp.Stdout)
	want := fmt.Sprintf("%d", req.TaskID)
	if got != want {
		t.Fatalf("list stdout=%q want %q", resp.Stdout, want)
	}
}
```
