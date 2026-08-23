# Scenario

**Issue**: normal completion remains unavailable from create stage

```
create -> tsk done <id> -> transition error
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "unfinished", "", nil)
	req.TaskID = id
	req.Args = []string{"done", fmt.Sprintf("%d", id)}
	return nil
}
```
