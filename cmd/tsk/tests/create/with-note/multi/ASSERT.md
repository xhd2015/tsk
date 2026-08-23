## Expected

- Two note lines; `first` appears before `second` in notes.jsonl.

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
	taskDir := taskAbs(req, inboxTaskRel(1, "create", req.Title))
	data, err := os.ReadFile(filepath.Join(taskDir, "notes.jsonl"))
	if err != nil {
		t.Fatalf("read notes.jsonl: %v", err)
	}
	s := string(data)
	i1 := strings.Index(s, "first note")
	i2 := strings.Index(s, "second note")
	if i1 < 0 || i2 < 0 || i1 > i2 {
		t.Fatalf("expected first before second in %q", s)
	}
}
```
