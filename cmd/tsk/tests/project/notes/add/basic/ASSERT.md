## Expected

- Exit 0; note text present; footer `1 notes`.
- On-disk `projects/1/notes.jsonl` exists.
- Auto or registry row has `"id": 1`.

## Exit Code

- 0

```go
import (
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
	assertContains(t, resp.Stdout, "dev command: go run ./")
	assertContains(t, resp.Stdout, "1 notes\n")

	notesPath := filepath.Join(req.TskHome, "projects", "1", "notes.jsonl")
	if _, err := os.Stat(notesPath); err != nil {
		t.Fatalf("missing %s: %v", notesPath, err)
	}
	reg, err := os.ReadFile(filepath.Join(req.TskHome, "projects.json"))
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, string(reg), `"id": 1`)
}
```
