# Scenario

**Feature**: listing notes on a task with no journal is empty, exit 0

```
create -> note list --id -> 0 notes
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "empty notes", "", nil)
	req.Args = []string{"note", "list", "--id", fmt.Sprintf("%d", id)}
	return nil
}
```
