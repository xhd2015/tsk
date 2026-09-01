# Scenario

**Feature**: `tsk tree` prints all tasks organized by topic tree

```
tsk tree [--json]
```

Inbox tasks (no topic) appear at the root level alongside top-level topics.

```go
import (
	"os"
	"os/exec"
	"strings"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	_ = initGitRepo
	return nil
}

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("output contains ANSI: %q", s)
	}
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
```

