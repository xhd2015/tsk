## Expected

- Exit code 0.
- Stdout lists `set` and `mkdir` subcommands.
- Stdout non-empty, ends with `\n`.
- Stderr empty.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "set")
	assertContains(t, resp.Stdout, "mkdir")
	assertContains(t, resp.Stdout, "where")
	assertContains(t, resp.Stdout, "info")
	assertContains(t, resp.Stdout, "notes")
	assertContains(t, resp.Stdout, "view")
}
```