# Scenario

**Feature**: `--note` does not match progress-labeled entries

```
progress add "blocked on review" -> search --note review -> 0 matches
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", fmt.Sprintf("%d", id), "--status", "blocked", "blocked on review")
	req.Args = []string{"search", "--note", "review"}
	return nil
}
```
