## Expected

- Exit code 0; stdout documents `--channel-id` and `--user`.
- Stdout does not document `--as`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "--channel-id")
	assertContains(t, resp.Stdout, "--user")
	assertNotContains(t, resp.Stdout, "--as")
}
```