## Expected

- Exit code 0; stdout `sent message 1\n`.
- `messages.jsonl` has one line: sender `alice`, body `fix login`, id `1`.
- `msg-counter` is `1`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "sent message 1")

	dir := activeChannelDir(req, "eng-alerts")
	msgs := readMessagesJSONL(t, dir)
	if len(msgs) != 1 {
		t.Fatalf("messages: got %d want 1", len(msgs))
	}
	if msgs[0].ID != 1 {
		t.Fatalf("message id: got %d want 1", msgs[0].ID)
	}
	if msgs[0].Sender != "alice" {
		t.Fatalf("sender: got %q want alice", msgs[0].Sender)
	}
	if msgs[0].Body != "fix login" {
		t.Fatalf("body: got %q want fix login", msgs[0].Body)
	}
	if readMsgCounter(t, dir) != 1 {
		t.Fatalf("msg-counter want 1")
	}
}
```
