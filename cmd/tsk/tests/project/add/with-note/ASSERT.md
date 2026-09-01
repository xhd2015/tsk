## Expected

- Exit 0; stdout id `1`.
- Task has one note with text `ctx`.

## Exit Code

- 0

```go
import (
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "1")
	taskDir := findTaskDirByID(t, req, 1)
	notesPath := filepath.Join(taskDir, "notes.jsonl")
	assertFileExists(t, notesPath)
	show := runTskOK(t, req, "show", "1")
	assertContains(t, show.Stdout, "notes: 1")
}
```
