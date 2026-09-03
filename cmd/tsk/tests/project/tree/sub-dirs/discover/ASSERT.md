## Expected

- Exit 0.
- `root-repo` and `nested-repo` project branches present with their tasks.
- `idle-repo` absent (no tasks).
- `deep-repo` absent (beyond default depth 3).
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
	assertNoANSI(t, resp.Stdout)
	assertContains(t, resp.Stdout, "root-repo  github.com/example/root-repo")
	assertContains(t, resp.Stdout, "nested-repo  github.com/example/nested-repo")
	assertContains(t, resp.Stdout, "from-root")
	assertContains(t, resp.Stdout, "from-nested")
	if strings.Contains(resp.Stdout, "idle-repo") {
		t.Fatalf("idle-repo should be omitted: %q", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "deep-repo") {
		t.Fatalf("deep-repo beyond depth 3 should be omitted: %q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "2 tasks, 2 projects")
}
```
