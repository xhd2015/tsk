## Expected

- Exit 0.
- `project: seatalk` (name fallback).
- No `cwd:` when task was created without recording a project add cwd (plain add).

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "project: seatalk")
	assertNotContains(t, resp.Stdout, "cwd:")
}
```
