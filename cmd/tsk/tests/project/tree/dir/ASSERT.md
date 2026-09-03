## Expected

- Exit 0; no ANSI.
- Tree shows `dot-pkgs` project and the `from-dir` task leaf.
- Footer `1 task, 1 project`.

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
	assertContains(t, resp.Stdout, "dot-pkgs  github.com/xhd2015/dot-pkgs")
	assertContains(t, resp.Stdout, "from-dir")
	assertContains(t, resp.Stdout, "1 task, 1 project")
}
```
