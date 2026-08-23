## Expected

- Exit code 0; stderr empty.
- JSON array one object: text `sess-abc`, labels `grok` and `session-id`; no ANSI.

## Exit Code

- 0

```go
import (
	"encoding/json"
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
	assertNoANSI(t, resp.Stdout)
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("json stdout should end with newline")
	}
	var arr []map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &arr); err != nil {
		t.Fatalf("parse notes json: %v; stdout=%q", err, resp.Stdout)
	}
	if len(arr) != 1 {
		t.Fatalf("len=%d stdout=%q", len(arr), resp.Stdout)
	}
	if arr[0]["text"] != "sess-abc" {
		t.Fatalf("text=%v", arr[0]["text"])
	}
	labels, _ := arr[0]["labels"].([]any)
	if len(labels) != 2 || labels[0] != "grok" || labels[1] != "session-id" {
		t.Fatalf("labels=%v", arr[0]["labels"])
	}
}
```
