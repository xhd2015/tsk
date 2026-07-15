## Expected

- Exit code 0; stdout `carol-room\n`.
- `participants.jsonl` has `carol` only (not alice).

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
	readChannelMetadata(t, dir)
	assertParticipantHandlesSorted(t, dir, []string{"carol"})
	for _, p := range readParticipantsJSONL(t, dir) {
		if p.Handle == "alice" {
			t.Fatalf("alice should not be participant when --user carol; got %+v", readParticipantsJSONL(t, dir))
		}
	}
}
```