# Scenario

**Feature**: `--json` emits structured events including `data.text` for notes

```
add; note add hello; logs --json -> array with mutation/action/data
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ExtraEnv = []string{"TSK_USER=alice"}
	id := addTask(t, req, "ship logs", "", nil)
	runTskOK(t, req, "note", "add", "--id", fmt.Sprintf("%d", id), "hello")
	req.TaskID = id
	req.Args = []string{"logs", "--json"}
	return nil
}
```
