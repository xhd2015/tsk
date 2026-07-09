## Expected

- Exit code 1.
- Stderr contains `task id required` exactly once (no duplicate from `fail()` + `main`).
- Stdout empty.

## Errors

- Single error line on stderr.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, stderr=%q", resp.Stderr)
	}
	if resp.Stdout != "" {
		t.Fatalf("stdout should be empty, got %q", resp.Stdout)
	}
	assertStderrContainsCount(t, resp.Stderr, "task id required", 1)
}
```