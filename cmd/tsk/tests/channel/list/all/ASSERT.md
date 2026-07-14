## Expected

- Exit code 0.
- Stdout contains both `active-one` and `archived-one`.
- Archived row shows `archived` status.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "active-one")
	assertContains(t, resp.Stdout, "archived-one")
	assertContains(t, resp.Stdout, "archived")
}
```
