## Expected

- Exit code 0; stdout `added bob\n`.
- `participants.jsonl` has `alice` and `bob` sorted.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "added bob")
	assertParticipantHandlesSorted(t, activeChannelDir(req, "team-ch"), []string{"alice", "bob"})
}
```
