## Expected

- Exit code 1.
- Stdout empty.
- Stderr is a single `Error: topic not found: no-such` line.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, stdout=%q", resp.Stdout)
	}
	if resp.Stdout != "" {
		t.Fatalf("stdout should be empty, got %q", resp.Stdout)
	}
	assertContains(t, resp.Stderr, "Error: topic not found: no-such")
	assertStderrContainsCount(t, resp.Stderr, "Error:", 1)
}
```
