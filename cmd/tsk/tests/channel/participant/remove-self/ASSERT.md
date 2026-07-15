## Expected

- Exit code 0; stdout `left team-ch\n`.
- Bob absent from `participants.jsonl`; alice remains.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "left team-ch")
	assertParticipantHandlesSorted(t, activeChannelDir(req, "team-ch"), []string{"alice"})
}
```