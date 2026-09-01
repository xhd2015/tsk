# Scenario

**Feature**: create under topic with one note (primary one-shot workflow)

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "flaky issue demo"
	req.Topic = "agent-pro"
	req.Args = []string{
		"add", "--topic", "agent-pro",
		"--note", "grok session abc track stall",
		"flaky issue demo",
	}
	return nil
}
```
