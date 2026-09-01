## Expected

- Non-zero exit; stderr mentions required flags.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit")
	}
	assertContains(t, resp.Stderr, "at least one of")
}
```
