## Expected

- Exit code 0.
- Every non-empty stdout line has rune width ≤ 36.
- `maxLineWidth(stdout)` ≤ 36.
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

	assertMaxWidth36(t, resp.Stdout)
	if w := maxLineWidth(resp.Stdout); w > 36 {
		t.Fatalf("maxLineWidth %d exceeds 36", w)
	}
}
```