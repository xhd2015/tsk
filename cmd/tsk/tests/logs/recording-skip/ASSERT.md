## Expected

- Exit 0.
- `events.jsonl` still has exactly one line (the prior `add`).

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	n := countEventLines(t, req)
	if n != 1 {
		t.Fatalf("events.jsonl lines=%d want 1 (logs must not append)", n)
	}
}
```
