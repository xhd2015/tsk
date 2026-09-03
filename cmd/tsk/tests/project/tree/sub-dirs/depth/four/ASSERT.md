## Expected

- Exit 0.
- `root-repo`, `nested-repo`, and `deep-repo` present.
- Footer `3 tasks, 3 projects`.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "root-repo  github.com/example/root-repo")
	assertContains(t, resp.Stdout, "nested-repo  github.com/example/nested-repo")
	assertContains(t, resp.Stdout, "deep-repo  github.com/example/deep-repo")
	assertContains(t, resp.Stdout, "from-deep")
	assertContains(t, resp.Stdout, "3 tasks, 3 projects")
}
```
