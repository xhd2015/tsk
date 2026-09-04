# Scenario

**Feature**: `note.add` is a mutation; `note.list` is not

```
add; note add hello; note list; logs -> add + note.add
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "journal", "", nil)
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "note", "add", "--id", idStr, "hello")
	runTskOK(t, req, "note", "list", "--id", idStr)
	req.TaskID = id
	req.Args = []string{"logs"}
	return nil
}
```
