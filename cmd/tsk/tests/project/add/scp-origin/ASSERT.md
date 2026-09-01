## Expected

- Exit 0.
- `project.origin` = `git.example.com/acme/loan-service/credit_backend/code-lens/tools/widget-cli`
- `project.name` absent

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "1")
	assertProjectOrigin(t, req, 1,
		"git.example.com/acme/loan-service/credit_backend/code-lens/tools/widget-cli",
		req.WorkRoot,
	)
}
```
