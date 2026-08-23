# Scenario

**Feature**: `tsk show` prints `notes:` count

```
create; two notes -> show -> notes: 2
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "show notes", "", nil)
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "note", "add", "--id", idStr, "one")
	runTskOK(t, req, "note", "add", "--id", idStr, "two")
	req.Args = []string{"show", idStr}
	return nil
}
```
