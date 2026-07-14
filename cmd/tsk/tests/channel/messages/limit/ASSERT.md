## Expected

- Exit code 0.
- Stdout contains `two`, not `one`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "two")
	assertNotContains(t, resp.Stdout, "one")
	assertContains(t, resp.Stdout, "[2]")
}
```
