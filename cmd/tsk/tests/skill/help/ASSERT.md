## Expected

- Exit 0; stderr empty.
- Usage mentions `--show`, `--install`, `--list`.
- Available topics includes `add` and `overview`.

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
	if !strings.HasPrefix(resp.Stdout, "Usage: tsk skill") {
		t.Fatalf("stdout=%q missing usage prefix", resp.Stdout)
	}
	for _, want := range []string{"--show", "--install", "--list", "Available topics:", "add", "overview"} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
}
```
