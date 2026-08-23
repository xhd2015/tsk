# Scenario

**Feature**: view of a missing topic errors

```
tsk topic view no-such
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"topic", "view", "no-such"}
	return nil
}
```
