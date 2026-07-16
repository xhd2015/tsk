## Expected

- Exit code 0; stdout `sent message 1\n`.
- Message sender is `bob`; body `from bob`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "sent message 1")
	msgs := readMessagesJSONL(t, activeChannelDir(req, "team-ch"))
	if len(msgs) != 1 {
		t.Fatalf("messages: got %d want 1", len(msgs))
	}
	if msgs[0].Sender != "bob" {
		t.Fatalf("sender: got %q want bob", msgs[0].Sender)
	}
	if msgs[0].Body != "from bob" {
		t.Fatalf("body: got %q want from bob", msgs[0].Body)
	}
}
```
