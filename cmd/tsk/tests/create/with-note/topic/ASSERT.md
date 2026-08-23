## Expected

- Exit 0; stdout id `1`.
- Task under `topics/agent-pro/`; `notes.jsonl` contains the note text.
- `show` reports `notes: 1`.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "1")

	wantRel := topicTaskRel("agent-pro", 1, "create", req.Title)
	taskDir := taskAbs(req, wantRel)
	assertDirExists(t, taskDir)
	assertIndexEquals(t, req, 1, wantRel)
	assertTopicPathEquals(t, req, 1, []string{"agent-pro"})

	data, err := os.ReadFile(filepath.Join(taskDir, "notes.jsonl"))
	if err != nil {
		t.Fatalf("read notes.jsonl: %v", err)
	}
	if !strings.Contains(string(data), "grok session abc track stall") {
		t.Fatalf("notes.jsonl missing text: %s", data)
	}
}
```
