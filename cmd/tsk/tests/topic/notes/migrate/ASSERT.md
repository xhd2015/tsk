## Expected

- Exit code 0; stderr empty.
- Stdout contains `legacy blob` and `2026-07-09T03:00:00Z`.
- `topic.json` no longer has a `notes` string.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertContains(t, resp.Stdout, "legacy blob")
	assertContains(t, resp.Stdout, "2026-07-09T03:00:00Z")
	assertContains(t, resp.Stdout, "1 notes\n")
	p := filepath.Join(req.TskHome, "topics", "knowledge-base", "topic.json")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read topic.json: %v", err)
	}
	var meta map[string]any
	if err := json.Unmarshal(data, &meta); err != nil {
		t.Fatalf("parse topic.json: %v", err)
	}
	if _, ok := meta["notes"]; ok {
		if s, _ := meta["notes"].(string); s != "" {
			t.Fatalf("legacy notes blob still set: %v", meta["notes"])
		}
	}
}
```
