# Scenario

**Feature**: channel with recent message activity stays active (no notify)

```
# last message created_at within idle threshold
check-channel-activity -> reads messages.jsonl -> status: active -> no exec
```

## Context

Leaves seed an active channel with a message timestamp within the default `1h` idle window.

```go
func Setup(t *testing.T, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```