## Expected

- Exit code 0; stderr empty.
- `stage: clarification`; `clarification[doing]`.
- `advance: blocked`.
- `next:` mentions `clarify confirm` (gated stage).
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
	assertAgentFact(t, resp.Stdout, "stage", "clarification")
	assertAgentFact(t, resp.Stdout, "terminal", "false")
	assertAgentDoing(t, resp.Stdout, "clarification")
	assertAgentAdvanceBlocked(t, resp.Stdout)
	assertAgentNextMentions(t, resp.Stdout, "clarify", "confirm")
}
```
