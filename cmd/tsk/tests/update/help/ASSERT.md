## Expected

- Exit 0; help mentions `--set-project` and `--set-topic`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "--set-project")
	assertContains(t, resp.Stdout, "--set-topic")
	assertContains(t, resp.Stdout, "--clear-topic")
}
```
