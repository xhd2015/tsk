# Scenario

**Feature**: done help documents forced completion

```
tsk done --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"done", "--help"}
	return nil
}
```
