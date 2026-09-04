## Expected

- Exit 0; `notes` group first under `@ seatalk`, then `[1] one`.
- Note text present; footer `1 task, 1 project`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	want := "" +
		".\n" +
		"└── @ seatalk\n" +
		"    ├── notes\n" +
		"    │   └── 2026-07-09T01:00:00Z  dev command: go run ./ --dev --port 8080\n" +
		"    └── [1] one  (create)\n" +
		"\n" +
		"1 task, 1 project\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
