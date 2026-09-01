# Scenario

**Feature**: `tsk project list` with empty registry

```
tsk project list -> 0 projects
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"project", "list"}
	return nil
}
```
