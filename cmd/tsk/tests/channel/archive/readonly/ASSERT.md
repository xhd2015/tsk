## Expected

- Exit code 1; stderr `Error:`.
- No new messages in archive dir.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	msgs := readMessagesJSONL(t, archiveChannelDir(req, "readonly-ch"))
	if len(msgs) != 0 {
		t.Fatalf("expected no messages, got %d", len(msgs))
	}
}
```
