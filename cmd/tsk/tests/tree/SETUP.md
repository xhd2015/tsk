# Scenario

**Feature**: `tsk tree` prints all tasks organized by topic tree

```
tsk tree [--json]
```

Inbox tasks (no topic) appear at the root level alongside top-level topics.

```go
import "strings"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("output contains ANSI: %q", s)
	}
}
```
