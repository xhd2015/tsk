## Expected

- Exit 0; lines start with `1.` and `2.`; footer `2 notes`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "1.  ")
	assertContains(t, resp.Stdout, "2.  ")
	assertContains(t, resp.Stdout, "first")
	assertContains(t, resp.Stdout, "second")
	assertContains(t, resp.Stdout, "2 notes\n")
}
```
