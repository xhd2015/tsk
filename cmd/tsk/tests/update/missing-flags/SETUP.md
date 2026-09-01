# Scenario

**Feature**: update without flags errors

```
tsk update 1 -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "x", "", nil)
	req.Args = []string{"update", "1"}
	return nil
}
```
