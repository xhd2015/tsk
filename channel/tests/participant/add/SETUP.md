# Scenario

**Feature**: add participant to channel

```
seed -> AddParticipant(bob) by alice -> roster alice,bob sorted
```

## Steps

1. Seed channel; add bob.

```go
func Setup(t *testing.T, req *Request) error {
	seedChannel(t, req, "Team", "team-ch")
	req.Op = "participant_add"
	req.ChannelID = "team-ch"
	req.Sender = "alice"
	req.Handle = "bob"
	return nil
}
```