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
		"├── [1]-create-alpha  task 1  create\n" +
		"├── [2]-create-beta  task 2  create\n" +
		"└── # knowledge-base  aliases: 知识库\n" +
		"    ├── [3]-create-report  task 3  create\n" +
		"    └── # reports\n" +
		"        └── [4]-create-draft  task 4  create\n" +
		"\n" +
		"4 tasks, 1 topic, 0 projects\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
