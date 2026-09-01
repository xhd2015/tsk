# Scenario

**Feature**: `--index` is required for note edit

```
tsk note edit --id 1 "text" -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	req.Args = []string{"note", "edit", "--id", "1", "some text"}
	return nil
}
```
