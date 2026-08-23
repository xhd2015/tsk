## Expected

- Exit code 0; stderr empty.
- Root `.` then `kb` topic then task leaf.
- Under task: `notes` with one entry, `progress` with one `(in-progress)` entry.
- Footer `1 task, 1 topic`.

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
		"└── kb\n" +
		"    └── [1]-create-report  task 1  create\n" +
		"        ├── notes\n" +
		"        │   └── 2026-07-09T02:00:00Z  session abc\n" +
		"        └── progress\n" +
		"            └── 2026-07-09T01:00:00Z  [progress]  (in-progress)  investigating\n" +
		"\n" +
		"1 task, 1 topic\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
