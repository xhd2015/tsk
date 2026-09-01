# Scenario

**Feature**: `tsk install` writes ~/.local/bin wrappers and ensures PATH

```
tsk install pmark
  -> ~/.local/bin/pmark (#!/bin/sh → tsk project add)
  -> bash/zsh rc PATH checker for ~/.local/bin
```

## Preconditions

- Leaves set `HOME` to `WorkRoot` via `ExtraEnv` so installs stay sandboxed.
- Shared helpers: `installHomeEnv`, `localBinPath`.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureInstallHelpersUsed()
	req.ExtraEnv = append(req.ExtraEnv, "HOME="+req.WorkRoot)
	return nil
}

func ensureInstallHelpersUsed() {
	_ = installHomeEnv
	_ = localBinPath
	_ = readLocalBin
	_ = installInitGitRepo
	_ = installAssertProjectOrigin
}

func installHomeEnv(req *Request) {
	has := false
	for _, e := range req.ExtraEnv {
		if strings.HasPrefix(e, "HOME=") {
			has = true
			break
		}
	}
	if !has {
		req.ExtraEnv = append(req.ExtraEnv, "HOME="+req.WorkRoot)
	}
}

func localBinPath(req *Request, name string) string {
	return filepath.Join(req.WorkRoot, ".local", "bin", name)
}

func readLocalBin(t *testing.T, req *Request, name string) string {
	t.Helper()
	body, err := os.ReadFile(localBinPath(req, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(body)
}

func installInitGitRepo(t *testing.T, dir, origin string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init")
	run("remote", "add", "origin", origin)
}

func installAssertProjectOrigin(t *testing.T, req *Request, id int, wantOrigin, wantCwdPrefix string) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	if task.Project == nil {
		t.Fatalf("task %d: project is nil", id)
	}
	if task.Project.Origin != wantOrigin {
		t.Fatalf("task %d project.origin: got %q want %q", id, task.Project.Origin, wantOrigin)
	}
	if wantCwdPrefix != "" && !strings.HasPrefix(task.Cwd, wantCwdPrefix) && task.Cwd != wantCwdPrefix {
		// cwd may be tilde-form; accept either absolute prefix or non-empty
		if !strings.Contains(task.Cwd, filepath.Base(wantCwdPrefix)) {
			t.Fatalf("task %d cwd=%q want prefix %q", id, task.Cwd, wantCwdPrefix)
		}
	}
}
```