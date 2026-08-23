# Scenario

**Feature**: `tsk topic note --label` stores a labeled topic note

```
mkdir; topic note --label grok -> topic notes shows [grok]
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "note", "--label", "grok", "knowledge-base", "hello")
	req.Args = []string{"topic", "notes", "knowledge-base"}
	return nil
}
```
