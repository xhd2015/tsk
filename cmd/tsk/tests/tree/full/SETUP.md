# Scenario

**Feature**: inbox + topic with aliases + nested subtopic rendered as tree

```
create inbox tasks + topic with alias + nested subtopic -> tree
```

## Steps

1. `tsk add "alpha"` (inbox).
2. `tsk add "beta"` (inbox).
3. `tsk topic mkdir knowledge-base/reports`.
4. `tsk add --topic knowledge-base "report"`.
5. `tsk add --topic knowledge-base/reports "draft"`.
6. `tsk topic alias add knowledge-base 知识库`.
7. `tsk tree`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "alpha", "", nil)
	addTask(t, req, "beta", "", nil)
	runTskOK(t, req, "topic", "mkdir", "knowledge-base/reports")
	addTask(t, req, "report", "knowledge-base", nil)
	addTask(t, req, "draft", "knowledge-base/reports", nil)
	runTskOK(t, req, "topic", "alias", "add", "knowledge-base", "知识库")
	req.Args = []string{"tree"}
	return nil
}
```
