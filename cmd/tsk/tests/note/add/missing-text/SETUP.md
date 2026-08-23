# Scenario

**Feature**: `tsk note add --id` without text errors

```
create -> note add --id N -> text required
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "note me", "", nil)
	req.Args = []string{"note", "add", "--id", fmt.Sprintf("%d", id)}
	return nil
}
```
