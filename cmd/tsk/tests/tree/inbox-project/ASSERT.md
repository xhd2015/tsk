## Expected

- Project group `@ dot-pkgs  github.com/…` with task 1.
- Ungrouped inbox task 2 at root.
- Footer includes `1 project`.

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
		"├── [2] solo  (create)\n" +
		"└── @ dot-pkgs  github.com/xhd2015/dot-pkgs\n" +
		"    └── [1] p  (create)\n" +
		"\n" +
		"2 tasks, 0 topics, 1 project\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
