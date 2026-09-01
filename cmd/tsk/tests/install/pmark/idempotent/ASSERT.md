## Expected

- Exit 0.
- Wrapper body unchanged.
- `.zshrc` has exactly one BEGIN checker marker.
- Stderr says PATH already includes `~/.local/bin`.

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
	body := readLocalBin(t, req, "pmark")
	want := "#!/bin/sh\nexec tsk project add \"$@\"\n"
	if body != want {
		t.Fatalf("body=%q want %q", body, want)
	}
	zshrc, err := os.ReadFile(filepath.Join(req.WorkRoot, ".zshrc"))
	if err != nil {
		t.Fatalf("read .zshrc: %v", err)
	}
	if n := strings.Count(string(zshrc), "# ----- BEGIN ~/.local/bin checker -----"); n != 1 {
		t.Fatalf("want 1 checker, got %d:\n%s", n, zshrc)
	}
	assertContains(t, resp.Stderr, "PATH already includes ~/.local/bin")
}
```
