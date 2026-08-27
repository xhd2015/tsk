## Expected

- Exit 0; stderr empty.
- Stdout is exactly `0 labels\n`.

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
	if resp.Stdout != "0 labels\n" {
		t.Fatalf("stdout=%q", resp.Stdout)
	}
}
```
