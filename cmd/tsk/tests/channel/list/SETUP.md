# Scenario

**Feature**: `tsk channel list` shows active channels; `--all` includes archived

```
tsk channel list [--json] [--all] -> table or JSON array; tombstoned never listed
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
