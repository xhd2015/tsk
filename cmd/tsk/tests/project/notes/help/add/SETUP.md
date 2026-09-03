# Scenario

**Feature**: `tsk project notes add -h`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"project", "notes", "add", "-h"}
	return nil
}
```
