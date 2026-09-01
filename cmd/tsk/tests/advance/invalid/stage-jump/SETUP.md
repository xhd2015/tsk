# Scenario

**Feature**: jumping create→implementation errors without mutation

```
create -> stage 1 implementation -> error; dir and stage unchanged
```

## Steps

1. `tsk add "add dark mode"`.
2. `tsk stage 1 implementation` (invalid jump).

```go
import (
	"fmt"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "add dark mode"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"stage", fmt.Sprintf("%d", id), "implementation"}
	return nil
}
```