## Expected

- Exit code 0; stderr empty.
- JSON array with one object `text=hello`; no ANSI; trailing newline.

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
	if arr[0]["text"] != "hello" {
		t.Fatalf("text=%v", arr[0]["text"])
	}
	if _, ok := arr[0]["ts"].(string); !ok || arr[0]["ts"] == "" {
		t.Fatalf("missing ts: %v", arr[0])
	}
}
```
