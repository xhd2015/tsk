## Expected

- Exit code 0; stderr empty.
- `stage: user_followup`; terminal false.
- Diagram: `user_followup[doing]` (exactly one doing); spine through `summary` bare; `(done)` future.
- `advance: ok` toward `clarification` (refine).
- `next:` mentions advance (refine path) and `done`.
- No ANSI; no rectangle chrome.

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
	assertAgentFact(t, resp.Stdout, "stage", "user_followup")
	assertAgentFact(t, resp.Stdout, "terminal", "false")
	assertAgentDoing(t, resp.Stdout, "user_followup")
	assertAgentFuture(t, resp.Stdout, "done")
	for _, stage := range []string{"create", "in_process", "clarification", "implementation", "verification", "summary"} {
		assertAgentPastBare(t, resp.Stdout, stage)
	}
	assertAgentAdvanceOK(t, resp.Stdout, "clarification")
	assertAgentNextMentions(t, resp.Stdout, "advance", "done")
	assertAgentArtHasBackLine(t, resp.Stdout)
}
```
