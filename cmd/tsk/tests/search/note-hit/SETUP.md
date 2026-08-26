# Scenario

**Feature**: default search finds a grok session id in a task note

```
create + note session id -> search id -> 1 note match
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "optimize git", "eng/backend", nil)
	sid := "01a01e87-2a6c-7ad2-9f48-0e0524256332"
	runTskOK(t, req, "note", "add", "--label", "grok", "--label", "session-id", "--id", fmt.Sprintf("%d", id), sid, "(git optimization session)")
	req.Args = []string{"search", sid}
	return nil
}
```
