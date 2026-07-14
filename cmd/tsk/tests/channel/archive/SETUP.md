# Scenario

**Feature**: `tsk channel archive` moves active channel to archive/ (readonly)

```
tsk channel archive --channel-id ID -> archive/<id>/; index archive/<id>
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
