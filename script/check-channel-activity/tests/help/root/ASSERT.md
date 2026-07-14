## Expected

- Exit code 0; stderr empty.
- Stdout contains `Usage:`, `--channel-id`, `--exec-if-idle-1h`, `--idle`, `--forever`, `--dry-run`.
- Help for `--exec-if-idle-1h` mentions shell command line and quoting spaces.
- Stdout ends with `\n`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	for _, flag := range []string{
		"Usage:",
		"--channel-id",
		"--exec-if-idle-1h",
		"--idle",
		"--forever",
		"--interval",
		"--dry-run",
		"--tsk-home",
	} {
		assertContains(t, resp.Stdout, flag)
	}
	assertContains(t, resp.Stdout, "shell command line")
	assertContains(t, resp.Stdout, "quote")
}
```