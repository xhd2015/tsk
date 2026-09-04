## Expected

- Exit 0; stderr empty.
- Stdout contains Usage, `--all`, `--limit`, `--json`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "Usage: tsk logs")
	assertContains(t, resp.Stdout, "--all")
	assertContains(t, resp.Stdout, "--limit")
	assertContains(t, resp.Stdout, "--json")
}
```
