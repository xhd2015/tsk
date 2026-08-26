## Expected

- Exit code 0.
- No ANSI escapes; contains plain hit text and `1 match`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	assertContains(t, resp.Stdout, "plain-token-xyz")
	assertContains(t, resp.Stdout, "1 match\n")
}
```
