# Scenario

**Feature**: `tsk channel delete` removes channel dir, index, writes tombstone

```
tsk channel delete --channel-id ID -> tombstones/<id>.json; blocks id reuse
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
