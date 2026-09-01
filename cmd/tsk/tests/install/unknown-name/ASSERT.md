## Expected

- Non-zero exit.
- Stderr contains `Error:` and `unknown name`.

## Exit Code

- non-zero

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertContains(t, resp.Stderr, "Error:")
	assertContains(t, resp.Stderr, "unknown name")
	assertContains(t, resp.Stderr, "pmark")
}
```
