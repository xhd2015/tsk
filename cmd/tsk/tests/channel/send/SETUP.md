# Scenario

**Feature**: `tsk channel send` appends message when sender is participant and channel active

```
participant --user HANDLE -> tsk channel send --channel-id ID <message...> -> messages.jsonl + msg-counter
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
