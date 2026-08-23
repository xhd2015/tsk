# Scenario

**Feature**: a second topic note appends instead of replacing

```
note hello; note world -> two jsonl lines
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic note knowledge-base hello`.
3. `tsk topic note knowledge-base world`.
4. `tsk topic notes knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "note", "knowledge-base", "hello")
	runTskOK(t, req, "topic", "note", "knowledge-base", "world")
	req.Args = []string{"topic", "notes", "knowledge-base"}
	return nil
}
```
