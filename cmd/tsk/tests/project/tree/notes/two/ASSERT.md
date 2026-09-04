## Expected

- Exit 0; both notes oldest first; then the task.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	want := "" +
		".\n" +
		"└── @ seatalk\n" +
		"    ├── notes\n" +
		"    │   ├── 2026-07-09T01:00:00Z  first remark\n" +
		"    │   └── 2026-07-09T02:00:00Z  second remark\n" +
		"    └── [1] one  (create)\n" +
		"\n" +
		"1 task, 1 project\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
