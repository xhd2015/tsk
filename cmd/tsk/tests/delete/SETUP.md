# Scenario

**Feature**: `tsk delete` permanently removes a task dir and index entry

```
# leaf: delete <id> -> deleted id\n; dir + index gone
# nested parent without --recursive -> error; with --recursive -> subtree gone
tsk delete [--recursive] <id>
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
