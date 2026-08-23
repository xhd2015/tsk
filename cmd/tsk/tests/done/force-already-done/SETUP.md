# Scenario

**Issue**: forced completion does not rewrite an already done task

```
create -> force done -> force done -> transition error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "finished", "", nil)
	runTskOK(t, req, "done", "--force", "1")
	req.TaskID = 1
	req.Args = []string{"done", "--force", "1"}
	return nil
}
```
