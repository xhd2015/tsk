# Scenario

**Feature**: `tsk progress add` / `tsk progress list` / `tsk progress edit` / `tsk progress archive` / `tsk progress show` record and display progress entries on a task

```
tsk progress add --status STATUS --id ID <text...>
tsk progress list [--json] [--limit N] [--status STATUS] [--show-index] --id ID
tsk progress edit --status STATUS --id ID --index N [text...]
tsk progress archive --id ID --index N
tsk progress show --id ID
```

Progress entries are notes labeled `progress` with an optional `status` field.

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
