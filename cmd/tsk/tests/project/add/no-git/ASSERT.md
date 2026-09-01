## Expected

- Non-zero exit.
- Stderr `Error:` and mentions registered / register help.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit outside git repo")
	}
	if !strings.Contains(resp.Stderr, "Error:") {
		t.Fatalf("stderr missing Error: prefix: %q", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "register") {
		t.Fatalf("stderr should mention register: %q", resp.Stderr)
	}
}
```
