# Scenario

**Feature**: `tsk label list -h` documents list usage

```
tsk label list -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"label", "list", "-h"}
	return nil
}
```
