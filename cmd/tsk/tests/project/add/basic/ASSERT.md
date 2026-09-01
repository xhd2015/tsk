## Expected

- Exit 0; stdout `1`.
- Inbox task with `project.origin` only (no name), non-empty `cwd`.

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
	if resp.Stderr != "" {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "1")

	wantRel := inboxTaskRel(1, "create", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, 1, wantRel)
	assertTopicPathNull(t, req, 1)
	assertProjectOrigin(t, req, 1, "github.com/xhd2015/dot-pkgs", req.WorkRoot)
	assertLastEventCommand(t, req, "project")

	autoPath := filepath.Join(req.TskHome, "projects-auto.json")
	assertFileExists(t, autoPath)
	data, err := os.ReadFile(autoPath)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	assertContains(t, s, `"origin": "github.com/xhd2015/dot-pkgs"`)
	assertContains(t, s, `"first_seen_at": "2026-07-09T12:00:00+08:00"`)
	assertContains(t, s, `"last_seen_at": "2026-07-09T12:00:00+08:00"`)
}
```

