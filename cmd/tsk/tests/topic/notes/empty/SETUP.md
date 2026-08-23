# Scenario

**Feature**: listing notes on a topic with no journal is empty, exit 0

```
mkdir -> topic notes -> 0 notes
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic notes knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	req.Args = []string{"topic", "notes", "knowledge-base"}
	return nil
}
```
