# Scenario

**Feature**: first idle check runs notify command and writes state

```
# old message, no prior state
check-channel-activity -> status: notified -> touch marker -> state.json
```

## Steps

1. Seed active channel with old message.
2. Run one-shot check.

```go
func Setup(t *testing.T, req *Request) error {
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "stale", CreatedAt: oldActivityTS,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	req.Args = defaultCheckArgs(req)
	return nil
}
```