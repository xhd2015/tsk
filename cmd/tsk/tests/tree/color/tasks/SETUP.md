# Scenario

**Feature**: full tree styles root and nested done task leaves only

```
create root + nested tasks -> force both done -> tree --color
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "root done", "", nil)
	addTask(t, req, "nested done", "projects/archive", nil)
	addTask(t, req, "active", "projects", nil)
	runTskOK(t, req, "done", "--force", "1")
	runTskOK(t, req, "done", "--force", "2")
	req.Args = []string{"tree", "--color"}
	return nil
}
```
