## Expected

- Exit 0.
- Stdout has `origin:`, `name:`, `cwd:`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "origin:  github.com/xhd2015/wrk")
	assertContains(t, resp.Stdout, "name:    wrk")
	assertContains(t, resp.Stdout, "cwd:")
}
```
