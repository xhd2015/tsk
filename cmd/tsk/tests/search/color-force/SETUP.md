# Scenario

**Feature**: `tsk search --color` emits ANSI even on a non-TTY pipe

```
note hit -> search --color query -> green/blue/gray/bold SGR present
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "color demo", "eng/backend", nil)
	runTskOK(t, req, "note", "add", "--label", "grok", "--id", fmt.Sprintf("%d", id), "color-token-abc")
	req.Args = []string{"search", "--color", "color-token-abc"}
	return nil
}
```
