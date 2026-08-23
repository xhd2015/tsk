## Expected

- Exit code 0; stderr empty.
- Documents `view` and `--json`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "view")
	assertContains(t, resp.Stdout, "--json")
}
```
