# Scenario

**Feature**: `tsk logs` shows events.jsonl activity (default: mutations)

```
tsk add / list / note … -> tsk logs [--all|--json|--limit] -> activity lines
```

## Preconditions

- Isolated `TSK_HOME`; `TSK_DATE=2026-07-09` pins `ts` to `2026-07-09T12:00:00+08:00`.
- `tsk logs` itself is not appended to `events.jsonl`.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func countEventLines(t *testing.T, req *Request) int {
	t.Helper()
	path := filepath.Join(req.TskHome, "events.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatalf("read events.jsonl: %v", err)
	}
	s := strings.TrimSpace(string(data))
	if s == "" {
		return 0
	}
	return len(strings.Split(s, "\n"))
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	_ = countEventLines
	return nil
}
```
