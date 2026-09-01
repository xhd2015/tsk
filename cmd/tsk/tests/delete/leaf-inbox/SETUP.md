# Scenario

**Feature**: delete inbox leaf removes dir and index

```
create "oops" -> tsk delete 1 -> deleted 1; inbox dir + index/1 gone
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "oops"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"delete", fmt.Sprintf("%d", id)}
	return nil
}
```
