# Scenario

**Feature**: `tsk topic info` works on a mkdir-only topic with no topic.json

```
topic mkdir knowledge-base -> topic info knowledge-base
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic info knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	req.Topic = "knowledge-base"
	req.Args = []string{"topic", "info", "knowledge-base"}
	return nil
}
```
