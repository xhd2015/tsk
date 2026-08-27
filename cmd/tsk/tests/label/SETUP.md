# Scenario

**Feature**: `tsk label` manage task labels and list label names in use

```
tsk label add|rm|list
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
