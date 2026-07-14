## Expected

- Exit code 0.
- Stdout contains `active-one`, not `archived-one`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "active-one")
	assertNotContains(t, resp.Stdout, "archived-one")
}
```
