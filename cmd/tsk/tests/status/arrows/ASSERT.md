## Expected

- Exit code 0.
- Stdout contains at least 6 `▼` downward-flow arrows.
- Refine ends at **left mid** of clarification: `►│ clarification`.
- Refine starts at **left mid** of user_followup: `└─refine`.
- No-followup ends at **right mid** of done: line contains done box and `◄`.
- `user_followup` appears before terminal `◉`; no orphan `user_followup` box after `◉`.
- Exact geometry sealed by `status/diagram-golden`.
- Stderr empty.

## Exit Code

- 0

```go
import (
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertContainsArrowDown(t, resp.Stdout, 6)
	// New geometry: refine enters clarification from the left (not old right-rail ◄──).
	assertContains(t, resp.Stdout, "►│ clarification")
	assertContains(t, resp.Stdout, "└─refine")
	// no-followup enters done from the right
	if !strings.Contains(resp.Stdout, "◄") {
		t.Fatalf("expected ◄ into done (no-followup rail) in stdout:\n%s", resp.Stdout)
	}
	assertFollowupBeforeTerminal(t, resp.Stdout)
}
```
