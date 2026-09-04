# Scenario

**Feature**: leaf `--dry-run` prints would-delete and keeps dir + index

```
create "oops" -> delete --dry-run 1
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "oops"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"delete", fmt.Sprintf("%d", id), "--dry-run"}
	return nil
}
```
