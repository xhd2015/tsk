# Scenario

**Feature**: new message advances last_activity and allows notify again when still idle

```
# notify on old message -> append newer (still idle) message -> notify again
```

## Steps

1. Seed channel with old message; run first check (notified).
2. Append second message with newer but still-stale timestamp.
3. Run second check via `req.Args`.

```go
const midActivityTS = "2026-07-13T12:00:00Z"

func Setup(t *testing.T, req *Request) error {
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "stale", CreatedAt: oldActivityTS,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	first := runCheckOK(t, req, defaultCheckArgs(req)...)
	if !strings.Contains(first.Stdout, "status: notified") {
		t.Fatalf("first run: expected notified, got %q", first.Stdout)
	}
	assertFileExists(t, req.MarkerPath)
	if err := os.Remove(req.MarkerPath); err != nil {
		t.Fatalf("remove marker: %v", err)
	}

	dir := filepath.Join(channelsRoot(req), "active", req.ChannelID)
	msgs = append(msgs, channelMessage{
		ID: 2, Sender: "alice", Body: "still stale", CreatedAt: midActivityTS,
	})
	if err := writeMessages(dir, msgs); err != nil {
		t.Fatalf("append message: %v", err)
	}
	req.LastActivity = midActivityTS
	req.Args = defaultCheckArgs(req)
	return nil
}
```