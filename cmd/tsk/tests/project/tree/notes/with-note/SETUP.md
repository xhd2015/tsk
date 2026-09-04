# Scenario

**Feature**: one project note prints as `notes` group before task leaves

```
register; add task; notes add; tree --project --plain
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	projectAddOK(t, req, "one", "--project", "seatalk")
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "dev command: go run ./ --dev --port 8080")
	req.Args = []string{"project", "tree", "--project", "seatalk", "--plain"}
	return nil
}
```
