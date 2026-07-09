## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Facts: `id: <n>`, `stage: create`, `terminal: false`.
- Spine row in fixed order: `create` → `in_process` → `clarification` → `implementation` → `verification` → `summary` → `done`, joined by `->`.
- Current mark: `create[doing]`; later spine stages use future form `(name)`.
- No rectangle chrome (`+---`, `│ stage │`, Unicode box corners).
- No ANSI escapes.
- Advance guidance present (`advance:`).

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
	assertAgentDoing(t, resp.Stdout, "create")
	for _, stage := range []string{
		"in_process", "clarification", "implementation",
		"verification", "summary", "done",
	} {
		assertAgentFuture(t, resp.Stdout, stage)
	}

	idStr := fmt.Sprintf("%d", req.TaskID)
	assertAgentFact(t, resp.Stdout, "id", idStr)
	assertAgentFact(t, resp.Stdout, "stage", "create")
	assertAgentFact(t, resp.Stdout, "terminal", "false")
	assertAgentHasFactKeys(t, resp.Stdout, "advance")
	assertAgentAdvanceOK(t, resp.Stdout, "in_process")
}
```
