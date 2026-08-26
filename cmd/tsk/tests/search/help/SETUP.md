# Scenario

**Feature**: `tsk search -h` prints usage and kind flags

```
tsk search -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"search", "-h"}
	return nil
}
```
