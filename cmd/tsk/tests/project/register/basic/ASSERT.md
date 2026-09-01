## Expected

- Exit 0; stdout `registered seatalk`.
- `project list` shows the name.
- Duplicate register errors.

## Exit Code

- 0

```go
import "strings"

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

	dup := runTskCmd(t, req, "project", "register", "--name", "seatalk", "--cwd", outsideGitDir(t))
	if dup.ExitCode == 0 {
		t.Fatal("duplicate register should fail")
	}
	if !strings.Contains(dup.Stderr, "already registered") {
		t.Fatalf("stderr=%q", dup.Stderr)
	}
}
```
