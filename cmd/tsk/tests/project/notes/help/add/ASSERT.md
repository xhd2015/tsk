## Expected

- Exit 0; documents `--project` / `--dir` / `--label`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "--project")
	assertContains(t, resp.Stdout, "--dir")
	assertContains(t, resp.Stdout, "--label")
}
```
