## Expected

- Exit 0; frontmatter `name: tsk/add`; body mentions `--parent`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "name: tsk/add") {
		t.Fatalf("expected name: tsk/add:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "--parent") {
		t.Fatalf("expected --parent in create topic:\n%s", resp.Stdout)
	}
}
```
