## Expected

- Exit code 0; stderr empty.
- Art includes `user_followup`, `refine`, and `questions` (back line under spine).
- Art does **not** label a `satisfied` edge (satisfied is next-only when applicable; at create it must not appear on art).
- Spine row still present (`->` join).
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
	assertAgentArtHasBackLine(t, resp.Stdout)
	assertAgentArtNoSatisfied(t, resp.Stdout)
	// structural: followup branch hangs off summary (questions edge after summary on art)
	art := agentArtText(resp.Stdout)
	assertContains(t, art, "user_followup")
	assertContains(t, art, "questions")
	assertContains(t, art, "refine")
}
```
