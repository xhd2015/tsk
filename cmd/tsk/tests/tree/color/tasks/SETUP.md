# Scenario

**Feature**: full tree styles root and nested done task leaves only

```
create root + nested tasks -> force both done -> tree --color
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "root done", "", nil)
	createTask(t, req, "nested done", "projects/archive", nil)
	createTask(t, req, "active", "projects", nil)
	runTskOK(t, req, "done", "--force", "1")
	runTskOK(t, req, "done", "--force", "2")
	req.Args = []string{"tree", "--color"}
	return nil
}
```
