## Expected

- Exit code 1; stderr `Error:`.
- `participants.jsonl` unchanged (`alice` only).

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	assertParticipantHandlesSorted(t, activeChannelDir(req, req.ChannelID), []string{"alice"})
}
```