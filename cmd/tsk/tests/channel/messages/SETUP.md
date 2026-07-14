# Scenario

**Feature**: `tsk channel messages` shows transcript for participants only

```
participant -> tsk channel messages --channel-id ID [--json] [--limit N]
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
