## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Output is **diagram/plain** (ASCII boxes, stage box line), not agent facts.
- No agent whole-line `id:` + `title:` fact block.
- No ANSI (plain path).

## Exit Code

- 0

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertAutoFormatOK(t, resp, err)
	assertNotAgentFactBlock(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	assertBoxLineForStage(t, resp.Stdout, "create")
	// plain prefers ASCII chrome
	if !strings.Contains(resp.Stdout, "+") && !strings.Contains(resp.Stdout, "|") {
		t.Fatalf("expected plain/ASCII diagram chrome in:\n%s", resp.Stdout)
	}
}
```
