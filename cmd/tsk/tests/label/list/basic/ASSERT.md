## Expected

- Exit 0; stderr empty.
- Stdout lists sorted names `grok-session-id`, `progress`, `report` and footer `3 labels`.
- Does not list full `grok-session-id=abc-1` tokens.

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
	want := "grok-session-id\nprogress\nreport\n3 labels\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
