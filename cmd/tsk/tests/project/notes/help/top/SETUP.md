# Scenario

**Feature**: `tsk project notes -h`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"project", "notes", "-h"}
	return nil
}
```
