# Scenario

**Feature**: `tsk channel create` allocates channel id and on-disk layout

```
TSK_USER=alice -> tsk channel create <name> [--channel-id ID] -> active/<id>/ + index + empty messages.jsonl
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
