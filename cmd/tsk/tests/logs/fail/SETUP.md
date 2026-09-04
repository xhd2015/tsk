# Scenario

**Feature**: failed mutations still appear as `fail`

```
add; done; done again -> logs shows ok then fail
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "finish me", "", nil)
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "done", idStr)
	resp := runTskCmd(t, req, "done", idStr)
	if resp.ExitCode == 0 {
		t.Fatalf("second done should fail")
	}
	req.TaskID = id
	req.Args = []string{"logs"}
	return nil
}
```
