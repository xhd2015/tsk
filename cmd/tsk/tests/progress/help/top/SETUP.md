# Scenario

**Feature**: `tsk progress -h` lists add, list, and show

```
tsk progress -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"progress", "-h"}
	return nil
}
```
