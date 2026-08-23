# Scenario

**Feature**: topic note appends notes.jsonl; info shows the count

```
mkdir -> topic note -> topic info notes: 1
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic note knowledge-base hello world`.
3. `tsk topic info knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "note", "knowledge-base", "hello", "world")
	req.Args = []string{"topic", "info", "knowledge-base"}
	return nil
}
```
