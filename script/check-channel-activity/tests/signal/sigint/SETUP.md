# Scenario

**Feature**: SIGINT prints stopped and exits 0

```
# active channel (no notify) + --forever -> SIGINT -> stopped
```

## Steps

1. Seed active channel with recent message.
2. Run `--forever --interval 100ms` and send SIGINT via `req.SIGINTStop`.

```go
func Setup(t *testing.T, req *Request) error {
	ts := recentActivityTS()
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "ping", CreatedAt: ts,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	req.SIGINTStop = true
	req.Args = defaultCheckArgs(req, "--forever", "--interval", "100ms")
	return nil
}
```