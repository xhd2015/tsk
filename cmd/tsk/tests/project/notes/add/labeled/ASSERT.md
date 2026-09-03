## Expected

- Exit 0; line includes `[run]` and text; footer `1 notes`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "[run]")
	assertContains(t, resp.Stdout, "go test ./...")
	assertContains(t, resp.Stdout, "1 notes\n")
}
```
