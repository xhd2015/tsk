# Scenario

**Feature**: `tsk skill -h` prints skill help and topic index

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "-h"}
	return nil
}
```
