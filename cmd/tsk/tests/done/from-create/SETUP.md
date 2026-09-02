# Scenario

**Feature**: plain done completes from create stage

```
create -> tsk done <id> -> stage done
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "finish from create"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"done", fmt.Sprintf("%d", id)}
	return nil
}
```
