## Expected

- Exit code 0; stderr empty.
- Stdout is the `knowledge-base` directory, not a `知识库` folder.

## Exit Code

- 0

```go
import "path/filepath"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	want := filepath.Join(req.TskHome, "topics", "knowledge-base") + "\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
