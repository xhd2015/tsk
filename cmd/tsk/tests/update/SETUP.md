# Scenario

**Feature**: `tsk update` changes project and/or topic on an existing task

```
tsk update <id> [--set-project REF|--clear-project] [--set-topic PATH|--clear-topic]
```

```go
import (
	"os"
	"os/exec"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	_ = initGitRepo
	_ = assertProjectOrigin
	_ = assertProjectName
	_ = assertProjectNil
	return nil
}

func initGitRepo(t *testing.T, dir, origin string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	cmd := exec.Command("git", "-C", dir, "init")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "-C", dir, "remote", "add", "origin", origin)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git remote add: %v\n%s", err, out)
	}
}

func assertProjectOrigin(t *testing.T, req *Request, id int, wantOrigin string) {
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

func assertProjectNil(t *testing.T, req *Request, id int) {
	t.Helper()
	task := readTaskJSON(t, findTaskDirByID(t, req, id))
	if task.Project != nil {
		t.Fatalf("task %d: project should be nil, got %+v", id, task.Project)
	}
}
```
