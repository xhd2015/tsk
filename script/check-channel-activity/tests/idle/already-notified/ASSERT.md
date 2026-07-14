## Expected

- Exit code 0; stderr empty.
- Stdout `status: already notified`.
- Notify marker not touched (still absent).

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertStdoutEndsWithNewline(t, resp.Stdout)
	assertStatusBlock(t, resp.Stdout, "already notified", req.LastActivity)
	assertFileNotExists(t, req.MarkerPath)
}
```