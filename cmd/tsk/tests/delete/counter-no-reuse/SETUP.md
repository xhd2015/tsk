# Scenario

**Feature**: after delete, next create still allocates next id (no reuse)

```
create -> delete 1 -> create -> id 2
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "gone", "", nil)
	runTskOK(t, req, "delete", fmt.Sprintf("%d", id))
	req.Title = "next"
	req.Args = []string{"add", req.Title}
	return nil
}
```
