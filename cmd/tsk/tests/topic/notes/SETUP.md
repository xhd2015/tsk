# Scenario

**Feature**: `tsk topic notes` lists the topic journal

```
tsk topic notes [--json] [--limit N] <topic>
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
