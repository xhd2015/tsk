## Expected

- Exit code 0; stderr empty.
- Root `.` with one leaf `[1]-create-solo  task 1  create`.
- Footer `1 task, 0 topics` (singular).

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
		".\n" +
		"└── [1]-create-solo  task 1  create\n" +
		"\n" +
		"1 task, 0 topics, 0 projects\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
