## Expected

- Exit 0.
- Stdout `registered <basename>`.
- List shows that name.

## Exit Code

- 0

```go
import "fmt"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, fmt.Sprintf("registered %s", req.Title))

	list := runTskOK(t, req, "project", "list", "--registered")
	assertContains(t, list.Stdout, req.Title)
}
```
