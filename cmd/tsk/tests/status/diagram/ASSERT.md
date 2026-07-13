## Expected Output

Compact pipeline at `clarification` with `│ clarification │` inside a 3-line box. Edge labels `refine` and `confirmed` present. Green ANSI on the clarification line. No mermaid `stateDiagram` output. Full geometry sealed by `status/diagram-golden` (exact stdout); this leaf only checks color + presence.

## Expected

- Exit code 0.
- Stdout non-empty, ends with `\n`.
- Compact box art (`╭`); not mermaid-wide lines.
- `│ clarification │` box line exists.
- Stdout contains stages: `create`, `in_process`, `clarification`, `implementation`, `verification`, `summary`, `user_followup`, `done`.
- Stdout contains edge labels `refine` and `confirmed`.
- Line containing `clarification` includes green ANSI `\x1b[32m`.
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

	assertCompactBoxArt(t, resp.Stdout)
	assertNotMermaidWide(t, resp.Stdout)
	assertBoxLineForStage(t, resp.Stdout, "clarification")
	assertStdoutHasStages(t, resp.Stdout,
		"create", "in_process", "clarification", "implementation",
		"verification", "summary", "user_followup", "done",
	)
	assertContains(t, resp.Stdout, "refine")
	assertContains(t, resp.Stdout, "confirmed")
	assertStageLineHasGreen(t, resp.Stdout, "clarification")
}
```
