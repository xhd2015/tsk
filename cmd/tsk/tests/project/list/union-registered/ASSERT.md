## Expected

- Default list shows registered name with `TASKS` column value `0`.
- `--auto` empty; `--registered` shows name without `TASKS`; `--all` same as default.
- `--all --auto` errors.

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
	assertContains(t, resp.Stdout, "TASKS")
	assertContains(t, resp.Stdout, "seatalk")
	assertContains(t, resp.Stdout, "1 project")
	lines := strings.Split(strings.TrimSpace(resp.Stdout), "\n")
	if len(lines) < 2 {
		t.Fatalf("stdout=%q", resp.Stdout)
	}
	if !strings.HasSuffix(strings.TrimRight(lines[1], " "), "0") {
		t.Fatalf("expected tasks 0: %q", lines[1])
	}

	auto := runTskOK(t, req, "project", "list", "--auto")
	assertStdoutTrimmedEquals(t, auto.Stdout, "0 projects")

	reg := runTskOK(t, req, "project", "list", "--registered")
	assertContains(t, reg.Stdout, "seatalk")
	assertContains(t, reg.Stdout, "NAME")
	assertNotContains(t, reg.Stdout, "TASKS")

	all := runTskOK(t, req, "project", "list", "--all")
	assertContains(t, all.Stdout, "seatalk")
	assertContains(t, all.Stdout, "TASKS")

	bad := runTskCmd(t, req, "project", "list", "--all", "--auto")
	if bad.ExitCode == 0 {
		t.Fatal("expected conflict")
	}
	if !strings.Contains(bad.Stderr, "mutually exclusive") {
		t.Fatalf("stderr=%q", bad.Stderr)
	}
}
```
