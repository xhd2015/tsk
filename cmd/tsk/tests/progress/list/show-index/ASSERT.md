## Expected

- Exit code 0; stderr empty.
- Entries are prefixed `1.  ` and `2.  `.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	want := "1.  2026-07-09T01:00:00Z  [progress]  (in-progress)  first\n" +
		"2.  2026-07-09T02:00:00Z  [progress]  (blocked)  second\n" +
		"2 entries\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
