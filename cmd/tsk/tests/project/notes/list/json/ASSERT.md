## Expected

- Exit 0; JSON array with text `hello`; no ANSI.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	assertContains(t, resp.Stdout, `"text":"hello"`)
	assertContains(t, resp.Stdout, `"ts":`)
}
```
