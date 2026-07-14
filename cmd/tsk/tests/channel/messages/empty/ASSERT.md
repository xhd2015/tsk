## Expected

- Exit code 0; stderr empty.
- Stdout does not contain message bodies.

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
	msgs := readMessagesJSONL(t, activeChannelDir(req, "quiet-ch"))
	if len(msgs) != 0 {
		t.Fatalf("expected no messages on disk, got %d", len(msgs))
	}
}
```
