## Expected

- Exit 0; stderr empty.
- Stdout documents `add` and `list`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "add")
	assertContains(t, resp.Stdout, "list")
}
```
