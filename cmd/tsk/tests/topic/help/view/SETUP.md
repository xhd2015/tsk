# Scenario

**Feature**: `tsk topic view --help`

```
tsk topic view --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"topic", "view", "--help"}
	return nil
}
```
