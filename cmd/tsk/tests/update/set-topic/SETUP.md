# Scenario

**Feature**: update --set-topic moves inbox task into topic

```
add inbox -> topic mkdir eng -> update --set-topic eng
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "into eng"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Topic = "eng"
	runTskOK(t, req, "topic", "mkdir", "eng")
	req.Args = []string{"update", fmt.Sprintf("%d", id), "--set-topic", "eng"}
	return nil
}
```
