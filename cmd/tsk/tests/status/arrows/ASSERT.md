## Expected

- Exit code 0.
- Stdout contains at least 6 `▼` downward-flow arrows.
- Stdout contains `►` or `──►` (summary→done branch).
- Stdout contains `◄──` (refine loop on clarification).
- `user_followup` appears before terminal `◉`; no orphan `user_followup` box after `◉`.
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

	assertContainsArrowDown(t, resp.Stdout, 6)
	if !strings.Contains(resp.Stdout, "►") && !strings.Contains(resp.Stdout, "──►") {
		t.Fatalf("expected ► or ──► branch arrow in stdout:\n%s", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "◄──")
	assertFollowupBeforeTerminal(t, resp.Stdout)
}
```