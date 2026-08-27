## Expected

- Exit code 0; stderr empty.
- Stdout contains `hit` and `1 notes`; does not contain `miss`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertContains(t, resp.Stdout, "hit")
	assertContains(t, resp.Stdout, "1 notes\n")
	assertNotContains(t, resp.Stdout, "miss")
}
```
