# Scenario

**Feature**: focused tree styles a done task leaf

```
create -> force done -> tree --id --color
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "finished", "projects", nil)
	runTskOK(t, req, "done", "--force", "1")
	req.Args = []string{"tree", "--id", "1", "--color"}
	return nil
}
```
