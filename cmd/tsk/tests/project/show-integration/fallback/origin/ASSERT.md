## Expected

- Exit 0.
- `project: github.com/xhd2015/unknown-proj`.
- No `cwd:` (plain add + update only).

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "project: github.com/xhd2015/unknown-proj")
	assertNotContains(t, resp.Stdout, "cwd:")
}
```
