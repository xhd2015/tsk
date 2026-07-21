## Expected

- Exit code 0.
- Stdout contains task id, title, stage `create`, label `bug`, and slug `show-me`.
- Stdout ends with `\n`; stderr empty.

## Exit Code

- 0

```go
import (
	"strings"
	"fmt"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("stdout should end with newline")
	}

	assertContains(t, resp.Stdout, fmt.Sprintf("%d", req.TaskID))
	assertContains(t, resp.Stdout, req.Title)
	assertContains(t, resp.Stdout, "create")
	assertContains(t, resp.Stdout, "bug")
	assertContains(t, resp.Stdout, "show-me")
}
```