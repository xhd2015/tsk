# Scenario

**Feature**: forced completion from create stage

```
create -> tsk done --force <id> -> *-done-*
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "finish directly"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"done", "--force", fmt.Sprintf("%d", id)}
	return nil
}
```
