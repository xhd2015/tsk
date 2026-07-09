## Expected

- Exit code 1.
- Stderr is a single error line mentioning format and/or allowed values (`diagram`, `agent`).
- Stdout empty or without a successful agent/diagram status dump (no success pipeline).

## Errors

- Invalid `--format` rejected once on stderr.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 1 {
		t.Fatalf("expected exit 1 for invalid --format, got %d stderr=%q stdout=%q",
			resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	if strings.TrimSpace(resp.Stdout) != "" {
		// allow empty; reject accidental success output
		if strings.Contains(resp.Stdout, "stage:") || strings.Contains(resp.Stdout, "[doing]") {
			t.Fatalf("stdout must not be a successful status dump, got %q", resp.Stdout)
		}
	}
	stderr := strings.TrimSpace(resp.Stderr)
	if stderr == "" {
		t.Fatal("stderr should explain invalid format")
	}
	// single stderr line (no duplicate fail+main)
	lines := 0
	for _, line := range strings.Split(stderr, "\n") {
		if strings.TrimSpace(line) != "" {
			lines++
		}
	}
	if lines != 1 {
		t.Fatalf("expected exactly one stderr line, got %d: %q", lines, resp.Stderr)
	}
	low := strings.ToLower(stderr)
	if !strings.Contains(low, "format") && !strings.Contains(low, "diagram") && !strings.Contains(low, "agent") {
		t.Fatalf("stderr should mention format or allowed values, got %q", resp.Stderr)
	}
	// preferred: list allowed values
	if !strings.Contains(low, "diagram") && !strings.Contains(low, "agent") {
		// still pass if "format" present; soft check already above
		_ = low
	}
}
```
