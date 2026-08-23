# Scenario

**Feature**: `tsk tree --id` on missing task errors

```
tsk tree --id 99 -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"tree", "--id", "99"}
	return nil
}
```
