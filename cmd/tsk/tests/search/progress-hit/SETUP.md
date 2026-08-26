# Scenario

**Feature**: `--progress` finds progress entry text

```
progress add "blocked on review" -> search --progress review -> 1 progress match
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", fmt.Sprintf("%d", id), "--status", "blocked", "blocked on review")
	req.Args = []string{"search", "--progress", "review"}
	return nil
}
```
