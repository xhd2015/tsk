# Scenario

**Feature**: `tsk add --parent` under a topic-scoped parent

```
create --topic kb "parent" -> create --parent 1 "child"
```

## Steps

```go
import (
	"fmt"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "kb")
	parentID := addTask(t, req, "parent report", "kb", nil)
	req.Args = []string{"add", "--parent", fmt.Sprintf("%d", parentID), "child detail"}
	return nil
}
```
