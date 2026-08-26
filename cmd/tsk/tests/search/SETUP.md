# Scenario

**Feature**: `tsk search` substring search across task titles, notes, progress, and topic notes

```
seed tasks/notes -> tsk search [kind flags] <query> -> matches on stdout
```

## Preconditions

- Inherits root `TSK_HOME` isolation and process-local `tsk` binary.
- Uses less-flags `Group(CollectParsedFlags(&kinds).Bool(--task/--note/--progress/--topic/--all))`.

## Context

Kind resolution: no kind flags or any `--all` → all surfaces; otherwise OR of listed kinds.
`--all` with other kind flags does not error; `--all` wins. Progress entries are notes
with label `progress` and are excluded from `--note`.

```go
import (
	"strings"
	"testing"
)

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.ContainsRune(s, '\x1b') {
		t.Fatalf("unexpected ANSI in %q", s)
	}
}
```
