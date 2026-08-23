# Scenario

**Feature**: `tsk note add` / `tsk note list` append and list labeled notes on a task

```
tsk note add [--label LABEL]... --id ID <text...>
tsk note list [--label LABEL]... --id ID
```

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
