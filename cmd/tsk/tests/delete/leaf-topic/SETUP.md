# Scenario

**Feature**: delete topic-placed leaf removes dir and index

```
create --topic eng/backend "x" -> delete 1 -> topics/… gone; index gone
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "x"
	req.Topic = "eng/backend"
	id := addTask(t, req, req.Title, req.Topic, nil)
	req.TaskID = id
	req.Args = []string{"delete", fmt.Sprintf("%d", id)}
	return nil
}
```
