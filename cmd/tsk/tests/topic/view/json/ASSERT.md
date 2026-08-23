## Expected

- Exit code 0; stderr empty.
- JSON object with `path`, empty `tasks`/`subtopics` arrays; no ANSI.

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
	var obj map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &obj); err != nil {
		t.Fatalf("parse view json: %v; stdout=%q", err, resp.Stdout)
	}
	if obj["path"] != "knowledge-base" {
		t.Fatalf("path=%v", obj["path"])
	}
	tasks, _ := obj["tasks"].([]any)
	subs, _ := obj["subtopics"].([]any)
	if len(tasks) != 0 || len(subs) != 0 {
		t.Fatalf("expected empty children: %v", obj)
	}
}
```
