## Expected

- Exit 0; stdout `updated <id>`.
- Task `project.origin` is `github.com/xhd2015/dot-pkgs`.

## Exit Code

- 0

```go
import "fmt"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, fmt.Sprintf("updated %d\n", req.TaskID))
	assertProjectOrigin(t, req, req.TaskID, "github.com/xhd2015/dot-pkgs")
}
```
