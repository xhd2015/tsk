# Scenario

**Feature**: `tsk project tree` discovers nested git projects under the scan root

```
root git + external/nested git with tasks
  -> default tree shows both; --no-sub-dirs / --sub-dirs-depth refine
```

## Preconditions

- Nested leaf fixtures call `setupNestedProjectTree`.
- Idle nested checkout (no tasks) must not appear as a project branch.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Sandbox scan_repo cache under WorkRoot (default is $HOME/.cache/git-repo-scan).
	req.ExtraEnv = append(req.ExtraEnv, "HOME="+req.WorkRoot)
	_ = setupNestedProjectTree
	return nil
}

// setupNestedProjectTree creates:
//   WorkRoot/                 <- origin root-repo (current)
//   WorkRoot/external/nested  <- origin nested-repo (has task)
//   WorkRoot/external/idle    <- origin idle-repo (no tasks; must stay hidden)
//   WorkRoot/a/b/c/deep       <- origin deep-repo (depth 4; for depth tests)
func setupNestedProjectTree(t *testing.T, req *Request) {
	t.Helper()
	initGitRepo(t, req.WorkRoot, "https://github.com/example/root-repo.git")

	nested := filepath.Join(req.WorkRoot, "external", "nested")
	idle := filepath.Join(req.WorkRoot, "external", "idle")
	deep := filepath.Join(req.WorkRoot, "a", "b", "c", "deep")
	initGitRepo(t, nested, "https://github.com/example/nested-repo.git")
	initGitRepo(t, idle, "https://github.com/example/idle-repo.git")
	initGitRepo(t, deep, "https://github.com/example/deep-repo.git")

	projectAddOK(t, req, "from-root")
	projectAddOK(t, req, "from-nested", "--dir", nested)
	projectAddOK(t, req, "from-deep", "--dir", deep)
}
```
