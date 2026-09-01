# Scenario

**Feature**: update --clear-topic moves task to inbox

```
add --topic eng -> update --clear-topic
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "move home"
	req.Topic = "eng"
	id := addTask(t, req, req.Title, req.Topic, nil)
	req.TaskID = id
	req.Args = []string{"update", fmt.Sprintf("%d", id), "--clear-topic"}
	return nil
}
```
