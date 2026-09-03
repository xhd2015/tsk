# Scenario

**Feature**: `tsk show` prints `project: <name>` when ledger has no location

```
register name-only with empty cwd/location -> update --set-project -> show
```

```go
import (
	"fmt"
	"os"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	// Strip path fields so show cannot resolve location.
	regPath := filepath.Join(req.TskHome, "projects.json")
	if err := os.WriteFile(regPath, []byte(`{
  "projects": [
    { "name": "seatalk" }
  ]
}
`), 0o600); err != nil {
		t.Fatal(err)
	}
	id := addTask(t, req, "name fallback", "", nil)
	runTskOK(t, req, "update", fmt.Sprintf("%d", id), "--set-project", "seatalk")
	req.TaskID = id
	req.Args = []string{"show", fmt.Sprintf("%d", id)}
	return nil
}
```
