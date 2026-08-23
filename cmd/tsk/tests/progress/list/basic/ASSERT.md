## Expected

- Exit code 0; stderr empty.
- Three entries with `(status)` in parens after `[progress]`.
- Footer `3 entries`.

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
	want := "2026-07-09T01:00:00Z  [progress]  (in-progress)  in-progress investigation\n" +
		"2026-07-09T02:00:00Z  [progress]  (in-progress)  optimized fetch\n" +
		"2026-07-09T03:00:00Z  [progress]  (blocked)  waiting on upstream\n" +
		"3 entries\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
