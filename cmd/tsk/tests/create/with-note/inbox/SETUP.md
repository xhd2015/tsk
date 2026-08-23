# Scenario

**Feature**: inbox create with --note

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "inbox with note"
	req.Args = []string{"create", "--note", "session pointer", "inbox with note"}
	return nil
}
```
