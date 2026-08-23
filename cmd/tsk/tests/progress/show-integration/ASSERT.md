## Expected

- Exit code 0; stderr empty.
- Stdout contains `notes: 1`.
- Stdout contains `progress: in-progress (1 entry)`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "notes: 1\n") {
		t.Fatalf("stdout=%q missing notes: 1", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "progress: in-progress (1 entry)\n") {
		t.Fatalf("stdout=%q missing progress line", resp.Stdout)
	}
}
```
