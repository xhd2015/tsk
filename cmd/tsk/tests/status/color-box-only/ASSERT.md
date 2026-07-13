## Expected

- Exit code 0.
- Stderr empty.
- Stdout ends with `\n`.
- Implementation stage has a box mid-row (`│ implementation │`).
- That mid-row includes green ANSI `\x1b[32m` on the current-stage box.
- The leading left refine-rail `│` on that mid-row is **not** inside the same SGR span as the box (whole-line coloring is incorrect).

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertBoxLineForStage(t, resp.Stdout, "implementation")
	assertBoxColoredLeftRailClear(t, resp.Stdout, "implementation")
}
```
