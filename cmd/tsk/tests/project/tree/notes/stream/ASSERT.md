## Expected

- Exit 0; `notes` under `@ dot-pkgs`; task leaf; footer.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "├── notes\n")
	assertContains(t, resp.Stdout, "from stream")
	assertContains(t, resp.Stdout, "[1] one")
	assertContains(t, resp.Stdout, "1 task, 1 project")
}
```
