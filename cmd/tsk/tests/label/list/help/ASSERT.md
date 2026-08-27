## Expected

- Exit 0; stderr empty.
- Stdout mentions `label list` and deduped / key behavior.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "label list")
	assertContains(t, resp.Stdout, "key=value")
}
```
