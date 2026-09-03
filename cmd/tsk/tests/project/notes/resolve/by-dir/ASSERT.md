## Expected

- Exit 0; note listed; auto ledger has id and origin.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "from-dir")
	assertContains(t, resp.Stdout, "1 notes\n")

	auto, err := os.ReadFile(filepath.Join(req.TskHome, "projects-auto.json"))
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, string(auto), `"id": 1`)
	assertContains(t, string(auto), `github.com/xhd2015/notes-demo`)

	if _, err := os.Stat(filepath.Join(req.TskHome, "projects", "1", "notes.jsonl")); err != nil {
		t.Fatalf("notes file: %v", err)
	}
}
```