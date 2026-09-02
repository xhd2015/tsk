# Scenario

**Issue**: archive refuses a done task (no terminal↔terminal)

```
create -> done -> archive -> already done
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "finished", "", nil)
	runTskOK(t, req, "done", fmt.Sprintf("%d", id))
	req.TaskID = id
	req.Args = []string{"archive", fmt.Sprintf("%d", id)}
	return nil
}
```
