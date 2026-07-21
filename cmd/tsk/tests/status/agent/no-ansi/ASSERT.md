## Expected

- Exit code 0; stderr empty.
- Stdout contains agent spine / facts (still agent mode).
- No ANSI escape sequences (`\x1b[`) despite `--color`.
- No rectangle chrome.

## Exit Code

- 0

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertAgentNoANSI(t, resp)
	assertAgentNoRectChrome(t, resp.Stdout)
	assertAgentSpineOrder(t, resp.Stdout)
	assertAgentDoing(t, resp.Stdout, "create")
	assertAgentFact(t, resp.Stdout, "stage", "create")
	if strings.Contains(resp.Stdout, "\x1b[") {
		t.Fatalf("agent format must ignore --color; found ANSI in:\n%s", resp.Stdout)
	}
}
```
