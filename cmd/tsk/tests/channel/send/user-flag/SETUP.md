# Scenario

**Feature**: `--user` overrides identity for send

```
add bob -> send --user bob "from bob" -> sender bob in jsonl
```

## Steps

1. Create, add bob, send with `--user bob`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	addParticipant(t, req, "team-ch", "bob")
	req.Args = []string{"channel", "send", "--channel-id", "team-ch", "--user", "bob", "from bob"}
	return nil
}
```