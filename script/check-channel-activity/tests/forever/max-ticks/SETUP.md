# Scenario

**Feature**: `--max-ticks` stops forever loop after N iterations

```
# tick1: notified -> tick2: already notified (test hook --max-ticks 2)
```

## Steps

1. Seed idle channel.
2. Run `--forever --interval 1ms --max-ticks 2`.

```go
func Setup(t *testing.T, req *Request) error {
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "stale", CreatedAt: oldActivityTS,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	req.Args = defaultCheckArgs(req,
		"--forever", "--interval", "1ms", "--max-ticks", "2",
	)
	return nil
}
```