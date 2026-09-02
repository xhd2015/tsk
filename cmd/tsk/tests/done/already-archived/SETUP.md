# Scenario

**Issue**: done refuses an archived (terminal) task

```
create -> archive -> done -> already archived
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "shelved", "", nil)
	runTskOK(t, req, "archive", fmt.Sprintf("%d", id))
	req.TaskID = id
	req.Args = []string{"done", fmt.Sprintf("%d", id)}
	return nil
}
```
