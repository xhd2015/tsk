## Expected

- Inbox task + one note line.

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
	wantRel := inboxTaskRel(1, "create", req.Title)
	taskDir := taskAbs(req, wantRel)
	assertDirExists(t, taskDir)
	data, err := os.ReadFile(filepath.Join(taskDir, "notes.jsonl"))
	if err != nil {
		t.Fatalf("read notes.jsonl: %v", err)
	}
	if !strings.Contains(string(data), "session pointer") {
		t.Fatalf("notes.jsonl missing text: %s", data)
	}
}
```
