# Scenario

**Feature**: state file prevents duplicate notify for same last_activity

```
# old message + matching state -> already notified -> no exec
```

## Steps

1. Seed active channel with old message.
2. Pre-write state with same `last_activity_at`.
3. Run one-shot check.

```go
func Setup(t *testing.T, req *Request) error {
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "stale", CreatedAt: oldActivityTS,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	writeNotifyState(t, req, req.LastActivity, "2026-07-14T10:00:00Z")
	req.Args = defaultCheckArgs(req)
	return nil
}
```