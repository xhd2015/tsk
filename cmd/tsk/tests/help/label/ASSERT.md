## Expected

- Exit code 0.
- Stdout lists `add` and `rm` subcommands.
- Stdout non-empty, ends with `\n`.
- Stderr empty.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "add")
	assertContains(t, resp.Stdout, "rm")
}
```