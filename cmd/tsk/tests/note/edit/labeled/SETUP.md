# Scenario

**Feature**: `tsk note edit --label` filters candidates, then `--index` selects

```
create -> note add (labeled) -> note add (unlabeled) -> edit by label+index -> list
```

## Steps

1. `tsk add "demo"`.
2. `tsk note add --id 1 --label grok "session abc"`.
3. `tsk note add --id 1 "other note"`.
4. `tsk note edit --id 1 --label grok --index 1 "session xyz"`.
5. `tsk note list --id 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "note", "add", "--id", "1", "--label", "grok", "session abc")
	runTskOK(t, req, "note", "add", "--id", "1", "other note")
	runTskOK(t, req, "note", "edit", "--id", "1", "--label", "grok", "--index", "1", "session xyz")
	req.Args = []string{"note", "list", "--id", "1"}
	return nil
}
```
