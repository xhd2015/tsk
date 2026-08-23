# Scenario

**Feature**: `tsk note add` appends a note to a task

```
create -> note add --id
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
