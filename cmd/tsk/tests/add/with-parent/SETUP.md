# Scenario

**Feature**: `tsk add --parent` nests under an inbox parent

```
create "parent" -> create --parent 1 "child" -> inbox/[1]-…/[2]-…
```

## Steps

1. Create inbox parent task.
2. Leaf Args create child with `--parent`.

```go
import (
	"fmt"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := addTask(t, req, "parent work", "", nil)
	req.Args = []string{"add", "--parent", fmt.Sprintf("%d", parentID), "child work"}
	return nil
}
```
