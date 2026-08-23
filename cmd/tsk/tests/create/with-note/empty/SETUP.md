# Scenario

**Feature**: empty --note text is rejected

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"create", "--note", "   ", "has empty note"}
	return nil
}
```
