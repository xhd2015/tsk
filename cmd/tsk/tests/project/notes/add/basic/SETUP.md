# Scenario

**Feature**: add then list shows the note; stored under `projects/<id>/notes.jsonl`

```
register --name seatalk; notes add --project seatalk "dev command: go run ./"; notes list
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "dev command: go run ./")
	req.Args = []string{"project", "notes", "list", "--project", "seatalk"}
	return nil
}
```
