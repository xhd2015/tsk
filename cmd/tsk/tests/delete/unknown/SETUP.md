# Scenario

**Feature**: delete unknown id fails

```
tsk delete 999 -> task 999 not found
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"delete", "999"}
	return nil
}
```
