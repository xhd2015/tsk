# Scenario

**Feature**: archive help documents any non-terminal completion

```
tsk archive --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"archive", "--help"}
	return nil
}
```
