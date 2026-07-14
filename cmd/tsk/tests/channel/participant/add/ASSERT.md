## Expected

- Exit code 0; stdout `added bob\n`.
- `channel.json` participants include `agent`, `alice`, `bob` sorted.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "added bob")
	ch := readChannelJSON(t, activeChannelDir(req, "team-ch"))
	assertChannelParticipantsSorted(t, ch, []string{"agent", "alice", "bob"})
}
```
