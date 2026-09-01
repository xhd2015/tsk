# Scenario

**Feature**: progress add requires --status

```
create task -> progress add without --status -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	req.Args = []string{"progress", "add", "--id", "1", "investigating"}
	return nil
}
```
