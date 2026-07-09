## Expected

- Exit code 0.
- Stdout equals top help (same content as `tsk` with no args): `Usage:`, full command list, command-specific help hint.
- Stdout non-empty, ends with `\n`.
- Stderr empty.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertTopHelpStdout(t, resp)
}
```