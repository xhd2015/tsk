# Scenario

**Feature**: `tsk topic where` prints the absolute topic directory

```
topic mkdir knowledge-base -> topic where knowledge-base
```

## Steps

1. `tsk topic mkdir knowledge-base`.
2. `tsk topic where knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	req.Topic = "knowledge-base"
	req.Args = []string{"topic", "where", "knowledge-base"}
	return nil
}
```
