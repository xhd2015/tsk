## Expected

- Exit code 1; stdout empty.
- Stderr contains `query required` exactly once with `Error:` prefix.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, stderr=%q", resp.Stderr)
	}
	if resp.Stdout != "" {
		t.Fatalf("stdout should be empty, got %q", resp.Stdout)
	}
	assertStderrContainsCount(t, resp.Stderr, "query required", 1)
	assertContains(t, resp.Stderr, "Error:")
}
```
