# Scenario

**Feature**: multiple --note flags preserve order

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "multi notes"
	req.Args = []string{
		"add",
		"--note", "first note",
		"--note", "second note",
		"multi notes",
	}
	return nil
}
```
