## Expected

- Exit code 0; stderr empty.
- Stdout mentions `task 1`, `note`, `eng/backend`, the session id, and `1 match`.

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
	assertContains(t, resp.Stdout, "task 1  note  eng/backend")
	assertContains(t, resp.Stdout, "labels=grok,session-id")
	assertContains(t, resp.Stdout, "01a01e87-2a6c-7ad2-9f48-0e0524256332")
	assertContains(t, resp.Stdout, "1 match\n")
}
```
