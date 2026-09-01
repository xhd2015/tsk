# Scenario

**Feature**: tree --id --color renders done and archived progress entries gray + strikethrough

```
create -> add done + archived -> tree --id --color
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "done", "completed")
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "archived", "retained history")
	req.Args = []string{"tree", "--id", "1", "--color"}
	return nil
}
```
