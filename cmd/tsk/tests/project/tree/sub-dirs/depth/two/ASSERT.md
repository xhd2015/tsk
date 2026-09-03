## Expected

- Exit 0.
- `root-repo` + `nested-repo`; `deep-repo` absent.
- Footer `2 tasks, 2 projects`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "nested-repo  github.com/example/nested-repo")
	if strings.Contains(resp.Stdout, "deep-repo") {
		t.Fatalf("deep-repo at depth 4 should be excluded at depth 2: %q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "2 tasks, 2 projects")
}
```
