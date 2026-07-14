## Expected

- Exit code 1.
- Stderr contains `Error:` (archived channel).
- No marker.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, stdout=%q stderr=%q", resp.Stdout, resp.Stderr)
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	assertFileNotExists(t, req.MarkerPath)
}
```