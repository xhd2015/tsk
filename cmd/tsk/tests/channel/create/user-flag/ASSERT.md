## Expected

- Exit code 0; stdout `carol-room\n`.
- `channel.json` participants are `agent` and `carol` sorted (not alice).

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
	assertStdoutTrimmedEquals(t, resp.Stdout, "carol-room")

	id := "carol-room"
	dir := activeChannelDir(req, id)
	assertDirExists(t, dir)
	ch := readChannelJSON(t, dir)
	assertChannelParticipantsSorted(t, ch, []string{"agent", "carol"})
	for _, p := range ch.Participants {
		if p.Handle == "alice" {
			t.Fatalf("alice should not be participant when --user carol; got %+v", ch.Participants)
		}
	}
}
```