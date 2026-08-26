## Expected

- Exit code 0.
- Stdout has `progress`, `status=blocked`, `blocked on review`, and `1 match`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "task 1  progress  inbox")
	assertContains(t, resp.Stdout, "status=blocked")
	assertContains(t, resp.Stdout, "blocked on review")
	assertContains(t, resp.Stdout, "1 match\n")
}
```
