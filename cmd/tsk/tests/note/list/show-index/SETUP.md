# Scenario

**Feature**: `tsk note list --show-index` prefixes 1-based index

```
create -> note add x2 -> note list --show-index
```

## Steps

1. `tsk create "demo"`.
2. `tsk note add --id 1 "first"`.
3. `tsk note add --id 1 "second"`.
4. `tsk note list --id 1 --show-index`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "note", "add", "--id", "1", "first")
	runTskOK(t, req, "note", "add", "--id", "1", "second")
	req.Args = []string{"note", "list", "--id", "1", "--show-index"}
	return nil
}
```
