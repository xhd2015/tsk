## Expected

- Exit code 0; stderr empty.
- Stdout is valid JSON array; no ANSI escape codes.
- Array contains channel with `id: json-room`.

## Exit Code

- 0

```go
import (
	"strings"
	"encoding/json"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("json stdout should end with newline, got %q", resp.Stdout)
	}
	var arr []map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &arr); err != nil {
		t.Fatalf("parse list json: %v; stdout=%q", err, resp.Stdout)
	}
	found := false
	for _, item := range arr {
		if id, _ := item["id"].(string); id == "json-room" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("json list missing json-room: %v", arr)
	}
}
```
