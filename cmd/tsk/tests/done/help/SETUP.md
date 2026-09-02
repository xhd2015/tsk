# Scenario

**Feature**: done help documents any non-terminal completion

```
tsk done --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"done", "--help"}
	return nil
}
```
