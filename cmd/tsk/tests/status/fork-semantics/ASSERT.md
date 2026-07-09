## Expected

- Exit code 0.
- Stdout contains `no followup` on the summary→done horizontal branch.
- `no followup` and `questions` do not appear on the same stdout line.
- The line containing `satisfied` also contains `►` (or `>` in `--plain`).
- Done box bottom does not use `╰──▼` (no ▼ embedded in box border).
- Stderr empty.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertForkSemantics(t, resp.Stdout)
}
```