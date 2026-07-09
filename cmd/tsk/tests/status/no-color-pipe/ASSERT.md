## Expected

- Exit code 0.
- No ANSI escape sequences (`\x1b[`).
- Compact box characters (`╭` or `│`).
- `│ clarification │` box line exists (current stage).
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
	assertCompactBoxArt(t, resp.Stdout)
	assertBoxLineForStage(t, resp.Stdout, "clarification")
	assertStdoutHasStages(t, resp.Stdout,
		"create", "in_process", "clarification", "implementation",
		"verification", "summary", "done",
	)
}
```