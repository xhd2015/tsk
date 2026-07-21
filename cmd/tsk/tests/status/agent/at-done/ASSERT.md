## Expected

- Exit code 0; stderr empty.
- `stage: done`; `terminal: true`.
- `done[doing]` on spine; other spine stages bare past.
- `advance: blocked` (or no advance-to); next empty or without further commands.
- No ANSI; no rectangle chrome.

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
	assertAgentFact(t, resp.Stdout, "stage", "done")
	assertAgentFact(t, resp.Stdout, "terminal", "true")
	assertAgentDoing(t, resp.Stdout, "done")
	for _, stage := range []string{
		"create", "in_process", "clarification", "implementation",
		"verification", "summary",
	} {
		assertAgentPastBare(t, resp.Stdout, stage)
	}
	assertAgentAdvanceBlocked(t, resp.Stdout)
	// next should not suggest advancing further
	plain := stripANSI(resp.Stdout)
	if strings.Contains(plain, "advance_to:") {
		t.Fatalf("done stage must not expose advance_to in:\n%s", resp.Stdout)
	}
}
```
