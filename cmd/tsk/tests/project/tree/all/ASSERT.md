## Expected

- Exit 0.
- Both `dot-pkgs` and `wrk` project branches appear.
- Footer `2 tasks, 2 projects`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "dot-pkgs  github.com/xhd2015/dot-pkgs")
	assertContains(t, resp.Stdout, "wrk  github.com/xhd2015/wrk")
	assertContains(t, resp.Stdout, "from-a")
	assertContains(t, resp.Stdout, "from-b")
	assertContains(t, resp.Stdout, "2 tasks, 2 projects")
}
```
