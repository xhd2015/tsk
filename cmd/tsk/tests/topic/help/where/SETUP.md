# Scenario

**Feature**: `tsk topic where --help` documents the command

```
tsk topic where --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"topic", "where", "--help"}
	return nil
}
```
