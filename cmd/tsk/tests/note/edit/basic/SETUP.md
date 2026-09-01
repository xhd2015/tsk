# Scenario

**Feature**: `tsk note edit` edits a note in place (replace text)

```
create -> note add -> note edit -> note list
```

## Steps

1. `tsk add "demo"`.
2. `tsk note add --id 1 "original text"`.
3. `tsk note edit --id 1 --index 1 "replaced text"`.
4. `tsk note list --id 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "note", "add", "--id", "1", "original text")
	runTskOK(t, req, "note", "edit", "--id", "1", "--index", "1", "replaced text")
	req.Args = []string{"note", "list", "--id", "1"}
	return nil
}
```
