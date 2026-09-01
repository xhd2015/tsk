# Scenario

**Feature**: topic rm errors when subtopics exist

```
mkdir eng/backend -> topic rm eng -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "eng/backend")
	req.Args = []string{"topic", "rm", "eng"}
	return nil
}
```
