## Expected

- Exit code 1.
- Stderr contains `Error:` once.
- Original channel dir unchanged.

## Errors

- Duplicate channel id rejected.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	if resp.Stdout != "" {
		t.Fatalf("stdout should be empty, got %q", resp.Stdout)
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	assertDirExists(t, activeChannelDir(req, "eng-alerts"))
}
```
