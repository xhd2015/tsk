## Expected

- Exit 1; stdout empty.
- Stderr mentions `Error:` and `empty key` once.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stdout != "" {
		t.Fatalf("stdout should be empty, got %q", resp.Stdout)
	}
	assertStderrContainsCount(t, resp.Stderr, "Error:", 1)
	assertContains(t, resp.Stderr, "empty key")
}
```
