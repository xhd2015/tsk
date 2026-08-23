# Scenario

**Feature**: `tsk topic info` counts tasks at the exact topic path

```
create --topic knowledge-base "x" -> topic info -> tasks: 1
```

## Steps

1. `tsk create --topic knowledge-base "x"`.
2. `tsk topic info knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "x"
	req.Topic = "knowledge-base"
	createTask(t, req, req.Title, req.Topic, nil)
	req.Args = []string{"topic", "info", "knowledge-base"}
	return nil
}
```
