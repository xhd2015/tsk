# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "--show", "working-on-task"}
	return nil
}
```
