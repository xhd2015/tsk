## Expected

- Exit code 0; no ANSI; valid JSON with kind note and text.

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
	assertNoANSI(t, resp.Stdout)
	var arr []map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &arr); err != nil {
		t.Fatalf("parse json: %v; stdout=%q", err, resp.Stdout)
	}
	if len(arr) != 1 || arr[0]["kind"] != "note" || arr[0]["text"] != "json-color-token" {
		t.Fatalf("arr=%v", arr)
	}
}
```
