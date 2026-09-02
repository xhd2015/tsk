# Scenario

**Feature**: archive from create stage

```
create -> tsk archive <id> -> stage archived
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "shelve me"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"archive", fmt.Sprintf("%d", id)}
	return nil
}
```
