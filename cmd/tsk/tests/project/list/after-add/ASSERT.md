## Expected

- Exit 0.
- Header includes `NAME`, `ORIGIN`, `LOCATION`, `TASKS`.
- Row includes `dot-pkgs`, `github.com/xhd2015/dot-pkgs`, and task count `1`.
- Footer `1 project`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "NAME")
	assertContains(t, resp.Stdout, "ORIGIN")
	assertContains(t, resp.Stdout, "LOCATION")
	assertContains(t, resp.Stdout, "TASKS")
	assertContains(t, resp.Stdout, "dot-pkgs")
	assertContains(t, resp.Stdout, "github.com/xhd2015/dot-pkgs")
	assertContains(t, resp.Stdout, "1 project")
	// TASKS column value for the single row
	lines := strings.Split(strings.TrimSpace(resp.Stdout), "\n")
	if len(lines) < 2 {
		t.Fatalf("stdout=%q", resp.Stdout)
	}
	data := strings.TrimRight(lines[1], " ")
	if !strings.HasSuffix(data, "1") {
		t.Fatalf("expected tasks=1 at end of data row: %q", lines[1])
	}

	active := runTskOK(t, req, "project", "list", "--active")
	assertContains(t, active.Stdout, "TASKS")
	assertContains(t, active.Stdout, "dot-pkgs")

	reg := runTskOK(t, req, "project", "list", "--registered")
	assertStdoutTrimmedEquals(t, reg.Stdout, "0 projects")

	autoOnly := runTskOK(t, req, "project", "list", "--auto")
	assertContains(t, autoOnly.Stdout, "dot-pkgs")
	assertContains(t, autoOnly.Stdout, "TASKS")
}
```
