## Expected

- Exit code 0; stderr empty.
- `notes: 1`.
- `notes.jsonl` exists and contains `hello world`.

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
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertContains(t, resp.Stdout, "notes: 1\n")
	p := filepath.Join(req.TskHome, "topics", "knowledge-base", "notes.jsonl")
	assertFileExists(t, p)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read notes.jsonl: %v", err)
	}
	if !strings.Contains(string(data), "hello world") {
		t.Fatalf("notes.jsonl missing text: %s", data)
	}
}
```
