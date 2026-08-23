## Expected

- Exit code 0; stderr empty.
- The single entry has preserved timestamp, `(done)`, and replacement text.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	want := "2026-07-09T01:00:00Z  [progress]  (done)  investigation complete\n1 entry\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
