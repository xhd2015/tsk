## Expected

- Exit code 0.
- Stdout `eng-alerts\n`.
- `channels/active/eng-alerts/channel.json` metadata only (`status: active`).
- `participants.jsonl` has creator `alice` only.
- `channels/index/eng-alerts` contains `active/eng-alerts`.
- `messages.jsonl` exists (empty); `msg-counter` exists.

## Side Effects

- Channel layout under `channels/active/eng-alerts/`.

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
	assertStdoutTrimmedEquals(t, resp.Stdout, "eng-alerts")

	id := "eng-alerts"
	dir := activeChannelDir(req, id)
	assertDirExists(t, dir)
	assertFileExists(t, filepath.Join(dir, "channel.json"))
	assertFileExists(t, filepath.Join(dir, "participants.jsonl"))
	assertFileExists(t, filepath.Join(dir, "messages.jsonl"))
	assertFileExists(t, filepath.Join(dir, "msg-counter"))
	assertChannelIndexEquals(t, req, id, "active/"+id)

	ch := readChannelMetadata(t, dir)
	if ch.ID != id {
		t.Fatalf("id: got %q want %q", ch.ID, id)
	}
	if ch.Name != req.ChannelName {
		t.Fatalf("name: got %q want %q", ch.Name, req.ChannelName)
	}
	if ch.Status != "active" {
		t.Fatalf("status: got %q want active", ch.Status)
	}
	assertParticipantHandlesSorted(t, dir, []string{"alice"})

	info, err := os.Stat(filepath.Join(dir, "messages.jsonl"))
	if err != nil {
		t.Fatalf("stat messages.jsonl: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("messages.jsonl should be empty, got %d bytes", info.Size())
	}
}
```