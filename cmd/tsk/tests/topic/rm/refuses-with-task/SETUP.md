# Scenario

**Feature**: topic rm errors when a task is under the topic

```
add --topic eng -> topic rm eng -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "stay", "eng", nil)
	req.Args = []string{"topic", "rm", "eng"}
	return nil
}
```
