# Scenario

**Feature**: update --help

```
tsk update --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"update", "--help"}
	return nil
}
```
