# Scenario

**Issue**: archive refuses an already archived task

```
create -> archive -> archive again -> already archived
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "twice", "", nil)
	runTskOK(t, req, "archive", fmt.Sprintf("%d", id))
	req.TaskID = id
	req.Args = []string{"archive", fmt.Sprintf("%d", id)}
	return nil
}
```
