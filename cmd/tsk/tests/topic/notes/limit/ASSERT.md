## Expected

- Exit code 0; stderr empty.
- Contains `two`; does not contain `one`.
- Footer `1 notes`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertContains(t, resp.Stdout, "two")
	assertNotContains(t, resp.Stdout, "one")
	assertContains(t, resp.Stdout, "1 notes\n")
}
```
