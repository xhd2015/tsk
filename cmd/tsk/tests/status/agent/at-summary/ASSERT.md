## Expected

- Exit code 0; stderr empty.
- `stage: summary`; `summary[doing]`.
- `advance: blocked` (gated: followup or done).
- `next:` includes followup path and done path.
- Art still has back line tokens `questions` / `user_followup` / `refine`.
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
	assertAgentFact(t, resp.Stdout, "stage", "summary")
	assertAgentFact(t, resp.Stdout, "terminal", "false")
	assertAgentDoing(t, resp.Stdout, "summary")
	assertAgentAdvanceBlocked(t, resp.Stdout)
	assertAgentNextMentions(t, resp.Stdout, "followup", "done")
	assertAgentArtHasBackLine(t, resp.Stdout)
}
```
