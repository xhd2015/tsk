## Expected

- Exit 0.
- Stdout contains `installed ~/.local/bin/pmark`.
- Wrapper exists, executable, body is `#!/bin/sh` + `exec tsk project add "$@"`.
- `.zshrc` contains the `~/.local/bin` checker begin marker.
- Stderr mentions PATH (added or already includes).

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
	assertContains(t, resp.Stdout, "installed ~/.local/bin/pmark")
	path := localBinPath(req, "pmark")
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat pmark: %v", err)
	}
	if st.Mode().Perm()&0o111 == 0 {
		t.Fatalf("pmark not executable: %v", st.Mode())
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
	if !strings.Contains(string(zshrc), "# ----- BEGIN ~/.local/bin checker -----") {
		t.Fatalf(".zshrc missing checker:\n%s", zshrc)
	}
	if !strings.Contains(resp.Stderr, "PATH") {
		t.Fatalf("stderr should mention PATH, got %q", resp.Stderr)
	}
}
```
