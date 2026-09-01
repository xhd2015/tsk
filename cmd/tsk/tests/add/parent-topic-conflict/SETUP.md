# Scenario

**Feature**: `--parent` and `--topic` together are rejected

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"add", "--parent", "1", "--topic", "kb", "x"}
	return nil
}
```
