## Expected

- Exit 1.
- Stderr: `Error:` and `--streaming conflicts with --json`.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stderr, "Error:")
	assertContains(t, resp.Stderr, "--streaming conflicts with --json")
}
```
