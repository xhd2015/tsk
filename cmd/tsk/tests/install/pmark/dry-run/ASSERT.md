## Expected

- Exit 0.
- Stdout has `[dry-run] would write ~/.local/bin/pmark` and would-create PATH checker lines.
- No `~/.local/bin/pmark` and no `.zshrc` on disk.

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
		t.Fatalf("exit %d stderr=%q stdout=%q", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	assertContains(t, resp.Stdout, "[dry-run] would write ~/.local/bin/pmark")
	assertContains(t, resp.Stdout, "[dry-run] would create ~/.zshrc (PATH checker)")
	if _, err := os.Stat(localBinPath(req, "pmark")); !os.IsNotExist(err) {
		t.Fatal("dry-run must not write pmark")
	}
	if _, err := os.Stat(filepath.Join(req.WorkRoot, ".zshrc")); !os.IsNotExist(err) {
		t.Fatal("dry-run must not create .zshrc")
	}
	if strings.Contains(resp.Stdout, "installed ") {
		t.Fatalf("live install line on dry-run: %q", resp.Stdout)
	}
}
```
