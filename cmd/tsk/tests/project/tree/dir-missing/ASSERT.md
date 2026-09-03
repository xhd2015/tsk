## Expected

- Exit code 1.
- Stderr mentions `resolve --dir` (or missing path).

## Errors

- Non-zero exit; error message on stderr.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, stderr=%q", resp.Stderr)
	}
	assertContains(t, resp.Stderr, "resolve --dir")
}
```
