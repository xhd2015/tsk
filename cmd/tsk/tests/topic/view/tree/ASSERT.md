## Expected

- Exit code 0; stderr empty.
- Root `knowledge-base`.
- Task `[1]-create-report-x` labeled `task 1  create`.
- Nested `reports` then `[2]-create-draft-slides` labeled `task 2  create`.
- Unicode tree branches; no `topic.json` / `notes.jsonl` names.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	want := "" +
		"knowledge-base\n" +
		"├── [1]-create-report-x  task 1  create\n" +
		"└── reports\n" +
		"    └── [2]-create-draft-slides  task 2  create\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
