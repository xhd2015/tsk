## Expected

- Topic `# eng` then `@ wrk  github.com/…/wrk` then task leaf.
- Footer `1 task, 1 topic, 1 project`.

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
		"└── # eng\n" +
		"    └── @ wrk  github.com/xhd2015/wrk\n" +
		"        └── [1]-create-report  task 1  create\n" +
		"\n" +
		"1 task, 1 topic, 1 project\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
