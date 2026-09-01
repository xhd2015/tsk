# Scenario

**Feature**: topic rm succeeds when topic has no tasks

```
topic mkdir eng -> topic rm eng
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "eng")
	req.Args = []string{"topic", "rm", "eng"}
	return nil
}
```
