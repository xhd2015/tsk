## Expected

- Exit code 0 (idempotent).
- Participants still exactly `alice`, `bob`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertParticipantHandlesSorted(t, activeChannelDir(req, req.ChannelID), []string{"alice", "bob"})
}
```