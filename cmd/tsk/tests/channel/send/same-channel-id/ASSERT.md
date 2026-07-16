## Expected

- Exit code 0; stdout `sent message 1\n`.
- One message with body `hi`, sender `alice`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "sent message 1")
	msgs := readMessagesJSONL(t, activeChannelDir(req, "eng-alerts"))
	if len(msgs) != 1 {
		t.Fatalf("messages: got %d want 1", len(msgs))
	}
	if msgs[0].Sender != "alice" {
		t.Fatalf("sender: got %q want alice", msgs[0].Sender)
	}
	if msgs[0].Body != "hi" {
		t.Fatalf("body: got %q want hi", msgs[0].Body)
	}
}
```
