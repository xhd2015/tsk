# Scenario

**Feature**: view shows tasks and nested sub-topics

```
create under knowledge-base + reports child task -> tree
```

## Steps

1. `tsk add --topic knowledge-base "report x"`.
2. mkdir `knowledge-base/reports`.
3. `tsk add --topic knowledge-base/reports "draft slides"`.
4. `tsk topic view knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "report x", "knowledge-base", nil)
	runTskOK(t, req, "topic", "mkdir", "knowledge-base/reports")
	addTask(t, req, "draft slides", "knowledge-base/reports", nil)
	req.Args = []string{"topic", "view", "knowledge-base"}
	return nil
}
```
