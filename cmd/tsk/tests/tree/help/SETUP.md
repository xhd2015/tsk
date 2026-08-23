# Scenario

**Feature**: `tsk tree -h` prints usage

```
tsk tree -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"tree", "-h"}
	return nil
}
```
