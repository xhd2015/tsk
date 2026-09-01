## Expected

- Exit 0.
- Stdout includes `cwd:` and `project: github.com/xhd2015/dot-pkgs` (origin only).

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "cwd:")
	assertContains(t, resp.Stdout, "project: github.com/xhd2015/dot-pkgs")
	assertContains(t, resp.Stdout, "show-me")
	assertNotContains(t, resp.Stdout, "(dot-pkgs)")
}
```
