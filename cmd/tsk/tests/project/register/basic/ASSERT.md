## Expected

- Exit 0; stdout `registered seatalk`.
- `project list` shows the name.
- Duplicate register with a different location errors (inconsistency).

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "registered seatalk")

	list := runTskOK(t, req, "project", "list")
	assertContains(t, list.Stdout, "NAME")
	assertContains(t, list.Stdout, "seatalk")
	assertContains(t, list.Stdout, "1 project")

	regPath := filepath.Join(req.TskHome, "projects.json")
	data, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, string(data), `"name": "seatalk"`)
	assertContains(t, string(data), `"location":`)
	assertContains(t, string(data), `"id": 1`)


	jsonList := runTskOK(t, req, "project", "list", "--registered", "--json")
	assertContains(t, jsonList.Stdout, `"location"`)

	dup := runTskCmd(t, req, "project", "register", "--name", "seatalk", "--cwd", outsideGitDir(t))
	if dup.ExitCode == 0 {
		t.Fatal("duplicate register should fail")
	}
	if !strings.Contains(dup.Stderr, "already registered") || !strings.Contains(dup.Stderr, "possible inconsistency") {
		t.Fatalf("stderr=%q", dup.Stderr)
	}
}
```
