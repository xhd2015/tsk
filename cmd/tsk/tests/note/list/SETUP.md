# Scenario

**Feature**: `tsk note list` reads a task journal

```
tsk note list --id ID
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
