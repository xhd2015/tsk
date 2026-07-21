# Scenario

**Feature**: move topic task back to inbox

```
create --topic eng/backend -> topic set 1 --inbox -> inbox/, topic_path null
```

## Steps

1. Create task under topic.
2. `tsk topic set 1 --inbox`.

```go
import (
	"fmt"
)

func Setup(t *testing.T, req *Request) error {
	req.Title = "return home"
	req.Topic = "eng/backend"
	id := createTask(t, req, req.Title, req.Topic, nil)
	req.TaskID = id
	req.Args = []string{"topic", "set", fmt.Sprintf("%d", id), "--inbox"}
	return nil
}
```