## Expected

- Exit code 0.
- Stderr empty.
- Stdout ends with `\n`.
- No ANSI escape sequences.
- Stdout is **byte-equal** to `expected.txt` (ASCII mapping of the unicode pipeline art; no-followup `+`/`|`/`+` column-aligned).
- `assertNoFollowupRailAligned`: no-followup corner, vertical right rail, and done-mid join share one column.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	assertStdoutEqualsFile(t, resp.Stdout, "expected.txt")
	assertNoFollowupRailAligned(t, resp.Stdout)
}
```