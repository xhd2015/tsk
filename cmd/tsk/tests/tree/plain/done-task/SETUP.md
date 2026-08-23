# Scenario

**Feature**: plain tree leaves a done task unstyled

```
create -> force done -> tree --plain
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "finished", "", nil)
	runTskOK(t, req, "done", "--force", "1")
	req.Args = []string{"tree", "--plain"}
	return nil
}
```
