# Scenario

**Feature**: `tsk logs` does not append to events.jsonl

```
add -> logs -> events.jsonl still one line (the add)
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "audit", "", nil)
	req.Args = []string{"logs"}
	return nil
}
```
