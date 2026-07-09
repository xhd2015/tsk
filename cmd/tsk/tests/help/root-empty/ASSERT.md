## Expected

- Exit code 0.
- Stdout contains `Usage:` and every root subcommand name.
- Stdout mentions `Run tsk <command> --help for command-specific options.`
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