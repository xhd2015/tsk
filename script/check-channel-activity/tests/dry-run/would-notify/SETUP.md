# Scenario

**Feature**: `--dry-run` on idle channel

```
check-channel-activity --dry-run -> status: would notify (dry-run)
```

## Steps

1. Seed idle channel.
2. Run with `--dry-run`.

```go
func Setup(t *testing.T, req *Request) error {
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "stale", CreatedAt: oldActivityTS,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	req.Args = defaultCheckArgs(req, "--dry-run")
	return nil
}
```