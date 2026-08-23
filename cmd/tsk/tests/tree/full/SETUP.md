# Scenario

**Feature**: inbox + topic with aliases + nested subtopic rendered as tree

```
create inbox tasks + topic with alias + nested subtopic -> tree
```

## Steps

1. `tsk create "alpha"` (inbox).
2. `tsk create "beta"` (inbox).
3. `tsk topic mkdir knowledge-base/reports`.
4. `tsk create --topic knowledge-base "report"`.
5. `tsk create --topic knowledge-base/reports "draft"`.
6. `tsk topic alias add knowledge-base 知识库`.
7. `tsk tree`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "alpha", "", nil)
	createTask(t, req, "beta", "", nil)
	runTskOK(t, req, "topic", "mkdir", "knowledge-base/reports")
	createTask(t, req, "report", "knowledge-base", nil)
	createTask(t, req, "draft", "knowledge-base/reports", nil)
	runTskOK(t, req, "topic", "alias", "add", "knowledge-base", "知识库")
	req.Args = []string{"tree"}
	return nil
}
```
