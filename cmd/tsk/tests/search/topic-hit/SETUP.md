# Scenario

**Feature**: `--topic` finds text in topic notes.jsonl

```
topic mkdir + note -> search --topic hello -> topic note match
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "note", "knowledge-base", "hello from topic journal")
	req.Args = []string{"search", "--topic", "hello"}
	return nil
}
```
