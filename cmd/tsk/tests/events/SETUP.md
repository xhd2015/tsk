# Scenario

**Feature**: every CLI invocation appends to events.jsonl

```
tsk <any> -> events.jsonl += one JSON line with command, args, exit_code
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```