## Expected

- Exit code 0; stderr empty.
- Root `.` with two inbox leaves then one topic node.
- Topic shows `aliases: 知识库`.
- Nested `reports` subtopic with one task leaf.
- Footer `4 tasks, 1 topic`.

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
		"├── [1] alpha  (create)\n" +
		"├── [2] beta  (create)\n" +
		"└── # knowledge-base  aliases: 知识库\n" +
		"    ├── [3] report  (create)\n" +
		"    └── # reports\n" +
		"        └── [4] draft  (create)\n" +
		"\n" +
		"4 tasks, 1 topic, 0 projects\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
