## Expected

- Exit code 0; stderr empty; no ANSI.
- JSON array one object: kind `note`, task_id 1, topic `eng/backend`, text `sess-json-1`.

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
		t.Fatalf("parse search json: %v; stdout=%q", err, resp.Stdout)
	}
	if len(arr) != 1 {
		t.Fatalf("len=%d stdout=%q", len(arr), resp.Stdout)
	}
	if arr[0]["kind"] != "note" {
		t.Fatalf("kind=%v", arr[0]["kind"])
	}
	if arr[0]["task_id"] != float64(1) {
		t.Fatalf("task_id=%v", arr[0]["task_id"])
	}
	if arr[0]["topic"] != "eng/backend" {
		t.Fatalf("topic=%v", arr[0]["topic"])
	}
	if arr[0]["text"] != "sess-json-1" {
		t.Fatalf("text=%v", arr[0]["text"])
	}
}
```
