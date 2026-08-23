## Expected

- Exit code 0; stderr empty.
- Stdout is a JSON array with one object containing `status`.
- No ANSI.

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
	assertNoANSI(t, resp.Stdout)
	want := "[{\"ts\":\"2026-07-09T01:00:00Z\",\"text\":\"investigating\",\"labels\":[\"progress\"],\"status\":\"in-progress\"}]\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
