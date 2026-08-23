## Expected

- Exit code 0; stderr empty.
- `subtopics: 1` and a `reports` child line.
- `tasks: 0` (exact path only).

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
	assertContains(t, resp.Stdout, "tasks: 0\n")
	assertContains(t, resp.Stdout, "subtopics: 1\n")
	assertContains(t, resp.Stdout, "  reports\n")
}
```
