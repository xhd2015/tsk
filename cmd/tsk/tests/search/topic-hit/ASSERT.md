## Expected

- Exit code 0.
- Stdout has `topic  note  knowledge-base` and the journal text; `1 match`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "topic  note  knowledge-base")
	assertContains(t, resp.Stdout, "hello from topic journal")
	assertContains(t, resp.Stdout, "1 match\n")
}
```
