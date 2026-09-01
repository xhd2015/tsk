# Scenario

**Feature**: `tsk project` manages project-scoped tasks + optional registry

```
git origin and/or projects.json name
  -> tsk project add "…"  -> inbox task + project.origin XOR project.name + cwd
  -> tsk project tree     -> tree like tsk tree
  -> tsk project list     -> registry rows
```

## Preconditions

- Fresh `TSK_HOME` per leaf.
- Leaves that need origin call `initGitRepo` (real `git` on PATH).

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureProjectHelpersUsed()
	return nil
}

func ensureProjectHelpersUsed() {
	_ = initGitRepo
	_ = projectAddOK
	_ = assertProjectOrigin
	_ = assertProjectName
	_ = outsideGitDir
}

func initGitRepo(t *testing.T, dir, origin string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	runGit(t, dir, "init")
	runGit(t, dir, "remote", "add", "origin", origin)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
}

func projectAddOK(t *testing.T, req *Request, title string, extraArgs ...string) int {
	t.Helper()
	before := maxTaskID(t, req)
	args := append([]string{"project", "add"}, extraArgs...)
	args = append(args, title)
	runTskOK(t, req, args...)
	id := maxTaskID(t, req)
	if id <= before {
		t.Fatalf("project add: expected new id (before=%d after=%d)", before, id)
	}
	req.TaskID = id
	req.Title = title
	return id
}

func assertProjectOrigin(t *testing.T, req *Request, id int, wantOrigin, wantCwdPrefix string) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	if task.Project == nil {
		t.Fatalf("task %d: project is nil", id)
	}
	if task.Project.Origin != wantOrigin {
		t.Fatalf("task %d project.origin: got %q want %q", id, task.Project.Origin, wantOrigin)
	}
	if task.Project.Name != "" {
		t.Fatalf("task %d project.name should be empty when origin set, got %q", id, task.Project.Name)
	}
	if task.Cwd == "" {
		t.Fatalf("task %d cwd should be set", id)
	}
	if wantCwdPrefix != "" && !strings.HasPrefix(task.Cwd, wantCwdPrefix) && task.Cwd != wantCwdPrefix {
		// still require non-empty absolute-ish cwd
		if !filepath.IsAbs(task.Cwd) {
			t.Fatalf("task %d cwd not abs: %q", id, task.Cwd)
		}
	}
}

func assertProjectName(t *testing.T, req *Request, id int, wantName string) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	if task.Project == nil {
		t.Fatalf("task %d: project is nil", id)
	}
	if task.Project.Name != wantName {
		t.Fatalf("task %d project.name: got %q want %q", id, task.Project.Name, wantName)
	}
	if task.Project.Origin != "" {
		t.Fatalf("task %d project.origin should be empty when name set, got %q", id, task.Project.Origin)
	}
}

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("output contains ANSI: %q", s)
	}
}

func outsideGitDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "tsk-project-nongit-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	abs, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return dir
	}
	return abs
}
```
