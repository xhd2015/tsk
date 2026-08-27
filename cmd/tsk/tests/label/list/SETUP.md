# Scenario

**Feature**: `tsk label list` prints deduped label names

```
tsk label list
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
