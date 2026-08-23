# Scenario

**Feature**: empty topic view prints `(empty)`

```
mkdir knowledge-base -> topic view
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic view knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	req.Args = []string{"topic", "view", "knowledge-base"}
	return nil
}
```
