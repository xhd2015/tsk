## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Stdout documents `--format` (and preferably values `diagram` / `agent`).

## Exit Code

- 0

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "--format")
	// values optional but strongly preferred
	low := strings.ToLower(resp.Stdout)
	if !strings.Contains(low, "agent") && !strings.Contains(low, "diagram") {
		t.Fatalf("status help should mention format values (agent/diagram), got:\n%s", resp.Stdout)
	}
}
```
