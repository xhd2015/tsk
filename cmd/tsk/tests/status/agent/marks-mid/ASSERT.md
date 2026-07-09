## Expected

- Exit code 0; stderr empty.
- `stage: implementation`; `implementation[doing]` exactly once.
- Past spine stages bare: `create`, `in_process`, `clarification` (not future, not doing).
- Future: `(verification)`, `(summary)`, `(done)`.
- Advance ok toward `verification`.
- No rectangle chrome; no ANSI.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertAgentNoANSI(t, resp)
	assertAgentNoRectChrome(t, resp.Stdout)
	assertAgentSpineOrder(t, resp.Stdout)
	assertAgentFact(t, resp.Stdout, "stage", "implementation")
	assertAgentFact(t, resp.Stdout, "terminal", "false")
	assertAgentDoing(t, resp.Stdout, "implementation")
	for _, stage := range []string{"create", "in_process", "clarification"} {
		assertAgentPastBare(t, resp.Stdout, stage)
	}
	for _, stage := range []string{"verification", "summary", "done"} {
		assertAgentFuture(t, resp.Stdout, stage)
	}
	assertAgentAdvanceOK(t, resp.Stdout, "verification")
	assertAgentNextMentions(t, resp.Stdout, "advance")
}
```
