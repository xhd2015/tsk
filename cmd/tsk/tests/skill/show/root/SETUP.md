# Scenario

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "--show"}
	return nil
}
```
