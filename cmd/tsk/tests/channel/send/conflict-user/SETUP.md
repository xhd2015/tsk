# Scenario

**Feature**: parent and leaf `--user` with different values → conflict error

```
# merge rule: set+set different user -> Error: conflicting …
tsk channel --channel-id team-ch --user alice send --user bob "msg"
  -> conflict error; no message written
```

## Steps

1. Create channel `team-ch`; add `bob` (alice already member as creator).
2. Run send with parent `--user alice` and leaf `--user bob` (both would be valid alone).

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	addParticipant(t, req, "team-ch", "bob")
	req.Args = []string{
		"channel", "--channel-id", "team-ch", "--user", "alice",
		"send", "--user", "bob", "msg",
	}
	return nil
}
```
