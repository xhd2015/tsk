## Expected

- Exit code 0.
- Finds the note `note-bar-body`; does not include progress text.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "note-bar-body")
	assertContains(t, resp.Stdout, "1 match\n")
	assertNotContains(t, resp.Stdout, "progress-baz-body")
	assertNotContains(t, resp.Stdout, "alpha-foo-title")
}
```
