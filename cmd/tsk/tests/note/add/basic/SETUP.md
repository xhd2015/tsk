# Scenario

**Feature**: unlabeled `tsk note add` then list shows the text

```
create -> note add --id N hello -> note list --id N
```

## Steps

1. `tsk create "note me"`.
2. `tsk note add --id <id> hello`.
3. `tsk note list --id <id>`.

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "note me", "", nil)
	runTskOK(t, req, "note", "add", "--id", fmt.Sprintf("%d", id), "hello")
	req.Args = []string{"note", "list", "--id", fmt.Sprintf("%d", id)}
	return nil
}
```
