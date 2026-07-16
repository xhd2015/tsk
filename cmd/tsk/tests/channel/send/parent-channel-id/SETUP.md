# Scenario

**Feature**: parent-level `--channel-id` before `send` is peeled and used for the channel

```
# parent peel: flags after channel, before subcommand
tsk channel --channel-id eng-alerts send "fix login"
  -> peel --channel-id -> run send with channel eng-alerts
  -> messages.jsonl sender alice, body "fix login"
```

## Steps

1. Create channel `eng-alerts` (alice creator/member).
2. Run `tsk channel --channel-id eng-alerts send "fix login"` (no leaf `--channel-id`).

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	req.Args = []string{"channel", "--channel-id", "eng-alerts", "send", "fix login"}
	return nil
}
```
