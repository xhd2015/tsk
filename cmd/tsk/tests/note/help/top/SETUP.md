# Scenario

**Feature**: `tsk note -h` lists add and list

```
tsk note -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"note", "-h"}
	return nil
}
```
