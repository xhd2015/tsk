## Expected

- Exit code 0.
- Every non-empty stdout line has rune width ≤ 42 (new geometry ~40; was 36).
- `maxLineWidth(stdout)` ≤ 42.
- Exact art is sealed by `status/diagram-golden` / `plain-golden`.
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

	assertMaxWidth42(t, resp.Stdout)
	if w := maxLineWidth(resp.Stdout); w > 42 {
		t.Fatalf("maxLineWidth %d exceeds 42", w)
	}
}
```
