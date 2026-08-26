## Expected

- Exit code 0.
- Stdout has `task 1  task  inbox` and title text; ends with `1 match`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "task 1  task  inbox")
	assertContains(t, resp.Stdout, "Optimize Git Clone")
	assertContains(t, resp.Stdout, "1 match\n")
}
```
