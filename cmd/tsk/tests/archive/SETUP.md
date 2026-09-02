# Scenario

**Feature**: `tsk archive` marks task terminal archived from any non-terminal stage

```
# any open stage -> archived; already terminal errors
tsk archive <id> -> stage archived
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
