## Expected

- Exit code 0.
- Stdout contains ASCII box characters (`+` or `|`).
- `| create |` box line exists (or tee-connected ASCII mid-row).
- No ANSI escape sequences (`\x1b[`).
- Soft structural check only; exact ASCII art sealed by `status/plain-golden`.
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

	assertNoANSI(t, resp.Stdout)
	assertContains(t, resp.Stdout, "+")
	assertBoxLineForStage(t, resp.Stdout, "create")
}
```
