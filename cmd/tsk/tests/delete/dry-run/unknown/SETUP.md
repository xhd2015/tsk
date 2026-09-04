# Scenario

**Feature**: `--dry-run` unknown id still fails

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"delete", "--dry-run", "999"}
	return nil
}
```
