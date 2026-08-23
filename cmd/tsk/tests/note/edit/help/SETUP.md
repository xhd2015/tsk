# Scenario

**Feature**: `tsk note edit -h` prints usage

```
tsk note edit -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"note", "edit", "-h"}
	return nil
}
```
