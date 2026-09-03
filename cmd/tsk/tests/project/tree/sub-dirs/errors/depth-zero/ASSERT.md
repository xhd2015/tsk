## Expected

- Exit 1.
- Stderr contains `Error:` and `--sub-dirs-depth must be >= 1`.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stderr, "Error:")
	assertContains(t, resp.Stderr, "--sub-dirs-depth must be >= 1")
}
```
