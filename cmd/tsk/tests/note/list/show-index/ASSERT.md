## Expected

- Exit code 0; stderr empty.
- Two notes prefixed with `1.  ` and `2.  `.
- Footer `2 notes`.

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
	want := "1.  2026-07-09T01:00:00Z  first\n" +
		"2.  2026-07-09T02:00:00Z  second\n" +
		"2 notes\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
