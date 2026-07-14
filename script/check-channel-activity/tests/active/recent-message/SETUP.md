# Scenario

**Feature**: recent message → active, no notify exec

```
# message created_at ~30m ago (< 1h threshold)
check-channel-activity -> status: active -> no marker touch
```

## Steps

1. Seed active channel with one recent message.
2. Run one-shot check.

```go
func Setup(t *testing.T, req *Request) error {
	ts := recentActivityTS()
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "ping", CreatedAt: ts,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	req.Args = defaultCheckArgs(req)
	return nil
}
```