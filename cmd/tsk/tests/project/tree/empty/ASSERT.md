## Expected

- Exit 0.
- Project node `agent-pro` present with `0 tasks, 1 project`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "agent-pro  github.com/xhd2015/agent-pro")
	assertContains(t, resp.Stdout, "0 tasks, 1 project")
}
```
