# Scenario

**Feature**: move inbox task into topic path

```
create (inbox) -> topic mkdir eng/backend -> topic set 1 eng/backend
```

## Steps

1. Create inbox task.
2. `tsk topic mkdir eng/backend`.
3. `tsk topic set 1 eng/backend`.

```go
import (
	"fmt"
)

func Setup(t *testing.T, req *Request) error {
	req.Title = "move me"
	topic := "eng/backend"
	req.Topic = topic
	id := createTask(t, req, req.Title, "", nil)
	req.Topic = topic
	runTskOK(t, req, "topic", "mkdir", topic)
	req.TaskID = id
	req.Args = []string{"topic", "set", fmt.Sprintf("%d", id), topic}
	return nil
}
```