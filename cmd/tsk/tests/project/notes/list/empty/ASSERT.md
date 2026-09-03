## Expected

- Exit 0; stdout exactly `0 notes\n`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stdout != "0 notes\n" {
		t.Fatalf("stdout=%q want %q", resp.Stdout, "0 notes\n")
	}
}
```
