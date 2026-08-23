# Scenario

**Feature**: `tsk topic` manages topic tree and task placement

```
# topic mkdir creates path; topic set moves task dir and updates topic_path + index
# where/info read the topic dir; note/alias write topics/<path>/topic.json
tsk topic mkdir <path>
tsk topic set <id> <path|--inbox>
tsk topic where <topic>
tsk topic info [--json] <topic>
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