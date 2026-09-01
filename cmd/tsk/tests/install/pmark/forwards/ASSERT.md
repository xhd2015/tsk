## Expected

- `tsk install pmark` exits 0 (Run).
- Executing the installed `pmark` prints a task id and stores `project.origin`.

## Exit Code

- 0

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("install exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	bin := getTskBin(t)
	pmark := localBinPath(req, "pmark")
	cmd := exec.Command(pmark, req.Title)
	cmd.Dir = req.WorkRoot
	path := filepath.Dir(bin) + string(os.PathListSeparator) +
		filepath.Join(req.WorkRoot, ".local", "bin") + string(os.PathListSeparator) +
		os.Getenv("PATH")
	env := tskEnv(req)
	filtered := make([]string, 0, len(env)+1)
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			continue
		}
		filtered = append(filtered, e)
	}
	filtered = append(filtered, "PATH="+path)
	cmd.Env = filtered
	out, err := captureCommandOutput(cmd)
	if err != nil {
		t.Fatalf("run pmark: %v", err)
	}
	if out.ExitCode != 0 {
		t.Fatalf("pmark exit %d stderr=%q stdout=%q", out.ExitCode, out.Stderr, out.Stdout)
	}
	idStr := strings.TrimSpace(out.Stdout)
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		t.Fatalf("pmark stdout id=%q stderr=%q", out.Stdout, out.Stderr)
	}
	installAssertProjectOrigin(t, req, id, "github.com/example/pmark-fwd", req.WorkRoot)
}
```
