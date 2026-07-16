# Scenario

**Feature**: parent-level `--channel-id` and `--user` before `send` set channel and sender

```
# parent peel both shared flags
tsk channel --channel-id team-ch --user bob send "from bob"
  -> peel --channel-id, --user -> send as bob
  -> messages.jsonl sender bob
```

## Steps

1. Create channel `team-ch`; add participant `bob`.
2. Run `tsk channel --channel-id team-ch --user bob send "from bob"` (no leaf flags).

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	addParticipant(t, req, "team-ch", "bob")
	req.Args = []string{"channel", "--channel-id", "team-ch", "--user", "bob", "send", "from bob"}
	return nil
}
```
