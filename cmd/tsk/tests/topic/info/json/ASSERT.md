## Expected

- Exit code 0; stderr empty.
- Stdout is one JSON object with `path`, `dir`, `tasks`; no ANSI; trailing newline.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"path/filepath"
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
		t.Fatalf("json stdout should end with newline, got %q", resp.Stdout)
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &obj); err != nil {
		t.Fatalf("parse info json: %v; stdout=%q", err, resp.Stdout)
	}
	if obj["path"] != "knowledge-base" {
		t.Fatalf("path: %v", obj["path"])
	}
	wantDir := filepath.Join(req.TskHome, "topics", "knowledge-base")
	if obj["dir"] != wantDir {
		t.Fatalf("dir: %v want %s", obj["dir"], wantDir)
	}
	tasks, ok := obj["tasks"].(float64)
	if !ok || tasks != 0 {
		t.Fatalf("tasks: %v", obj["tasks"])
	}
}
```
