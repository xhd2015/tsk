## Expected

- Exit 0.
- Only `root-repo` / `from-root`.
- Nested / deep / idle absent.
- Footer `1 task, 1 project`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "root-repo  github.com/example/root-repo")
	assertContains(t, resp.Stdout, "from-root")
	if strings.Contains(resp.Stdout, "nested-repo") {
		t.Fatalf("nested-repo should be hidden with --no-sub-dirs: %q", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "deep-repo") {
		t.Fatalf("deep-repo should be hidden with --no-sub-dirs: %q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "1 task, 1 project")
}
```
