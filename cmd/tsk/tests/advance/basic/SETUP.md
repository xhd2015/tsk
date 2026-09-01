# Scenario

**Feature**: advance from create updates stage in task.json (dirname unchanged)

```
create "add dark mode" -> advance 1 -> stage in_process; dir stays inbox/[1]-add-dark-mode/
```

## Steps

1. `tsk add "add dark mode"`.
2. `tsk advance 1`.

```go
import (
	"fmt"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "add dark mode"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"advance", fmt.Sprintf("%d", id)}
	return nil
}
```