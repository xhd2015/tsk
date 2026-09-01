# Scenario

**Feature**: `tsk search --no-color` suppresses ANSI

```
note hit -> search --no-color query -> no escapes
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "plain demo", "", nil)
	runTskOK(t, req, "note", "add", "--id", fmt.Sprintf("%d", id), "plain-token-xyz")
	req.Args = []string{"search", "--no-color", "plain-token-xyz"}
	return nil
}
```
