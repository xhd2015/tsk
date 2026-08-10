# Scenario

**Feature**: `tsk channel participant` manages channel roster

```
tsk channel participant add|remove --channel-id ID [handle]
tsk channel participants --channel-id ID [--json]
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
