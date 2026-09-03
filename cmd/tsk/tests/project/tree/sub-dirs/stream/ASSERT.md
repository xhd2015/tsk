## Expected

- Exit 0.
- Root project appears before nested.
- Idle/deep omitted; footer `2 tasks, 2 projects`.

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
	rootIdx := strings.Index(resp.Stdout, "root-repo  github.com/example/root-repo")
	nestedIdx := strings.Index(resp.Stdout, "nested-repo  github.com/example/nested-repo")
	if rootIdx < 0 || nestedIdx < 0 || rootIdx > nestedIdx {
		t.Fatalf("want root before nested: root=%d nested=%d out=%q", rootIdx, nestedIdx, resp.Stdout)
	}
	assertContains(t, resp.Stdout, "from-root")
	assertContains(t, resp.Stdout, "from-nested")
	if strings.Contains(resp.Stdout, "idle-repo") || strings.Contains(resp.Stdout, "deep-repo") {
		t.Fatalf("unexpected idle/deep in stream output: %q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "2 tasks, 2 projects")
}
```
