# Scenario

**Feature**: non-participant cannot add members

```
seed -> AddParticipant(bob) by carol (not a member) -> error
```

## Steps

1. Seed channel; carol attempts add.

```go
func Setup(t *testing.T, req *Request) error {
	seedChannel(t, req, "Private", "private-ch")
	req.Op = "participant_not_member"
	req.ChannelID = "private-ch"
	req.Sender = "carol"
	req.Handle = "bob"
	return nil
}
```