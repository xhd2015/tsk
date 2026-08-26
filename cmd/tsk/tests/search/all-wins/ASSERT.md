## Expected

- Exit code 0; stderr empty (no combine error).
- Finds the note despite `--task` because `--all` wins.

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
	assertContains(t, resp.Stdout, "unique-all-wins-token")
	assertContains(t, resp.Stdout, "note")
	assertContains(t, resp.Stdout, "1 match\n")
}
```
