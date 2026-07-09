## Expected Output

Pipeline includes workflow stage names in order. Current stage `clarification` is marked (e.g. `>` prefix, `*` suffix, or bracketed).

## Expected

- Exit code 0.
- Stdout non-empty, ends with `\n`.
- Stdout contains all major stage names: `create`, `in_process`, `clarification`, `implementation`, `verification`, `summary`, `done`.
- Stdout marks `clarification` as the current stage (implementation-defined marker adjacent to that stage name).
- Stderr empty.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	if resp.Stdout == "" {
		t.Fatal("stdout should not be empty")
	}
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("stdout should end with newline, got %q", resp.Stdout)
	}

	stages := []string{"create", "in_process", "clarification", "implementation", "verification", "summary", "done"}
	for _, stage := range stages {
		assertContains(t, resp.Stdout, stage)
	}

	// Current-stage marker: line containing clarification should differ from other stage lines.
	lines := strings.Split(strings.TrimSuffix(resp.Stdout, "\n"), "\n")
	var clarificationLines []string
	for _, line := range lines {
		if strings.Contains(line, "clarification") {
			clarificationLines = append(clarificationLines, line)
		}
	}
	if len(clarificationLines) == 0 {
		t.Fatal("stdout should contain clarification stage line")
	}
	marked := false
	markers := []string{">", "*", "[", "(current)", "◀", "▶"}
	for _, line := range clarificationLines {
		for _, m := range markers {
			if strings.Contains(line, m) {
				marked = true
				break
			}
		}
	}
	if !marked {
		t.Fatalf("clarification stage should be marked in stdout:\n%s", resp.Stdout)
	}
}
```