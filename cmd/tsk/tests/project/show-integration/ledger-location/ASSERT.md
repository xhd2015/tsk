## Expected

- Exit 0.
- Stdout includes task `cwd:` (absolute WorkRoot) and `project:` equal to ledger location (WorkRoot; tilde if under `$HOME`).
- Does not print origin on the `project:` line.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"

	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	wantCwd := "cwd: " + pathfmt.TildeHome(filepath.Clean(req.WorkRoot))
	assertContains(t, resp.Stdout, wantCwd)
	wantProject := "project: " + pathfmt.TildeHome(filepath.Clean(req.WorkRoot))
	assertContains(t, resp.Stdout, wantProject)
	assertContains(t, resp.Stdout, "show-me")
	assertNotContains(t, resp.Stdout, "project: github.com/xhd2015/dot-pkgs")
	assertNotContains(t, resp.Stdout, "(dot-pkgs)")

	autoPath := filepath.Join(req.TskHome, "projects-auto.json")
	data, err := os.ReadFile(autoPath)
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, string(data), `"location":`)
}
```
