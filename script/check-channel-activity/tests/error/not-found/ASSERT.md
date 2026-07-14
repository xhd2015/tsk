## Expected

- Exit code 1.
- Stderr contains `Error:`.
- Stdout empty.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, stdout=%q stderr=%q", resp.Stdout, resp.Stderr)
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	if resp.Stdout != "" {
		t.Fatalf("stdout should be empty, got %q", resp.Stdout)
	}
	assertFileNotExists(t, req.MarkerPath)
}
```