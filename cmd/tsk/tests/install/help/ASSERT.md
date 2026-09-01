## Expected

- Exit 0; stderr empty.
- Stdout documents `pmark` and `tsk project add`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "pmark")
	assertContains(t, resp.Stdout, "tsk project add")
	assertContains(t, resp.Stdout, "~/.local/bin")
	assertContains(t, resp.Stdout, "--dry-run")
}
```
