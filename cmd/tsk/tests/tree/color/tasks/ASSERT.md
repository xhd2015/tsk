## Expected

- Exit code 0; stderr empty.
- Root and nested done task leaves are gray + struck through with no trailing stage text.
- Branch art and active task leaves are unstyled (no trailing stage either).

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}

	ansi, reset := "\x1b[90m\x1b[9m", "\x1b[0m"
	for _, text := range []string{
		"[1] root done",
		"[2] nested done",
	} {
		if !strings.Contains(resp.Stdout, ansi+text+reset) {
			t.Fatalf("stdout=%q missing styled task leaf %q", resp.Stdout, text)
		}
	}
	if strings.Contains(resp.Stdout, ansi+"├──") || strings.Contains(resp.Stdout, ansi+"└──") {
		t.Fatalf("branch art must not be styled: %q", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, ansi+"[3] active") {
		t.Fatalf("active task must not be styled: %q", resp.Stdout)
	}
}
```
