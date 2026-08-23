## Expected

- Exit code 0; stderr empty.
- Branch art is uncolored.
- Done and archived entry content is prefixed with gray + strikethrough and reset.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	ansi := "\x1b[90m\x1b[9m"
	reset := "\x1b[0m"
	if strings.Contains(resp.Stdout, ansi+"└──") {
		t.Fatalf("branch art must not be styled: %q", resp.Stdout)
	}
	for _, text := range []string{"(done)  completed", "(archived)  retained history"} {
		if !strings.Contains(resp.Stdout, ansi) || !strings.Contains(resp.Stdout, text+reset) {
			t.Fatalf("stdout=%q missing styled %q", resp.Stdout, text)
		}
	}
}
```
