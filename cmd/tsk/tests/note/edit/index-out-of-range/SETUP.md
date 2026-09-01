# Scenario

**Feature**: index out of range errors

```
create -> note add -> note edit --index 5 -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "note", "add", "--id", "1", "only note")
	req.Args = []string{"note", "edit", "--id", "1", "--index", "5", "text"}
	return nil
}
```
