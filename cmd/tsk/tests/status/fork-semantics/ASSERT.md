## Expected

- Exit code 0.
- Stdout contains `no followup` on the summary→done horizontal branch.
- Stdout contains `refine` on the user_followup→clarification left rail.
- `no followup` and `questions` do not appear on the same stdout line.
- `satisfied` is a **vertical** spine label (like `claim`); **no** `satisfied►` sideways decoration.
- Done box bottom does not use `╰──▼` (no ▼ embedded in box border).
- Terminal `◉` is a dead end (no refine/clarification content after it).
- Exact geometry sealed by `status/diagram-golden`.
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
