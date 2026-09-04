# Scenario

**Feature**: default `tsk logs` shows add, hides list

```
tsk add; tsk list; tsk logs -> one add mutation line
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "ship logs", "", nil)
	runTskOK(t, req, "list")
	req.TaskID = id
	req.Args = []string{"logs"}
	return nil
}
```
