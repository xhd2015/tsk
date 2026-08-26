## Expected

- Exit code 0; stderr empty.
- Stdout contains green (`\x1b[32m`), blue (`\x1b[34m`), gray (`\x1b[90m`), bold (`\x1b[1m`), and reset (`\x1b[0m`).
- Bold wraps the query substring `color-token-abc`.

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
	out := resp.Stdout
	for _, seq := range []string{"\x1b[32m", "\x1b[34m", "\x1b[90m", "\x1b[1m", "\x1b[0m"} {
		if !strings.Contains(out, seq) {
			t.Fatalf("missing SGR %q in %q", seq, out)
		}
	}
	if !strings.Contains(out, "\x1b[1mcolor-token-abc\x1b[0m") {
		t.Fatalf("expected bold query highlight, stdout=%q", out)
	}
	assertContains(t, out, "1 match")
}
```
