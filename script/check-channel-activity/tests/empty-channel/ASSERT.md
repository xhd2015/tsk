## Expected

- Exit code 0; stderr empty.
- `last_activity` equals channel `created_at` (no messages).
- `status: notified`; marker created.

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
	assertStatusBlock(t, resp.Stdout, "notified", oldCreatedAtTS)
	assertFileExists(t, req.MarkerPath)
}
```