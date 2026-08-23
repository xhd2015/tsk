## Expected

- Exit code 0; stderr empty.
- `path: knowledge-base`, `title: knowledge-base`, empty aliases/notes, `tasks: 0`, `subtopics: 0`.
- `dir:` is the absolute topic directory.
- No `topic.json` is created.

## Exit Code

- 0

```go
import (
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
	assertContains(t, resp.Stdout, "path: knowledge-base\n")
	assertContains(t, resp.Stdout, "title: knowledge-base\n")
	assertContains(t, resp.Stdout, "aliases:\n")
	assertContains(t, resp.Stdout, "notes: 0\n")
	assertContains(t, resp.Stdout, "tasks: 0\n")
	assertContains(t, resp.Stdout, "subtopics: 0\n")
	wantDir := filepath.Join(req.TskHome, "topics", "knowledge-base")
	assertContains(t, resp.Stdout, "dir: "+wantDir+"\n")
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("stdout should end with newline")
	}
	assertFileNotExists(t, filepath.Join(wantDir, "topic.json"))
}
```
