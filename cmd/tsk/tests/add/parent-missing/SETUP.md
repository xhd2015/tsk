# Scenario

**Feature**: `create --parent` rejects missing parent

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"add", "--parent", "99", "orphan"}
	return nil
}
```
