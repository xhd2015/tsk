# Scenario

**Feature**: `tsk topic note` appends notes.jsonl

```
tsk topic note <topic> <text...>
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
