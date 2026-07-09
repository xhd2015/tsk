## Expected

- Exit code 0.
- Stdout contains `create` usage text.
- Stdout documents `--label` and `--topic` flags.
- Stdout non-empty, ends with `\n`.
- Stderr empty.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "create")
	assertContains(t, resp.Stdout, "--label")
	assertContains(t, resp.Stdout, "--topic")
}
```