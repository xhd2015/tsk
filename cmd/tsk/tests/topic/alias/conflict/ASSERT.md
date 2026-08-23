## Expected

- Exit code 1.
- Stdout empty.
- Stderr mentions `Error:` and that the alias is already used.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit")
	}
	if resp.Stdout != "" {
		t.Fatalf("stdout should be empty, got %q", resp.Stdout)
	}
	assertContains(t, resp.Stderr, "Error:")
	assertContains(t, resp.Stderr, "知识库")
}
```
