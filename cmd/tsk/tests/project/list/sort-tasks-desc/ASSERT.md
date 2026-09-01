## Expected

- Exit 0.
- Data rows: `alpha` (2 tasks) before `beta` (1 task).

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	lines := strings.Split(strings.TrimSpace(resp.Stdout), "\n")
	// header + 2 data + footer
	if len(lines) < 4 {
		t.Fatalf("stdout=%q", resp.Stdout)
	}
	if !strings.Contains(lines[1], "alpha") {
		t.Fatalf("first data row want alpha: %q", lines[1])
	}
	if !strings.HasSuffix(strings.TrimRight(lines[1], " "), "2") {
		t.Fatalf("alpha tasks want 2: %q", lines[1])
	}
	if !strings.Contains(lines[2], "beta") {
		t.Fatalf("second data row want beta: %q", lines[2])
	}
	if !strings.HasSuffix(strings.TrimRight(lines[2], " "), "1") {
		t.Fatalf("beta tasks want 1: %q", lines[2])
	}
}
```
