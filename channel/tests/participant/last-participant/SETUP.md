# Scenario

**Feature**: cannot remove the last participant

```
seed (alice only) -> RemoveParticipant(alice) -> error
```

## Steps

1. Seed channel with sole creator; attempt self-remove.

```go
func Setup(t *testing.T, req *Request) error {
	seedChannel(t, req, "Solo", "solo-ch")
	req.Op = "participant_remove"
	req.ChannelID = "solo-ch"
	req.Sender = "alice"
	req.Handle = "alice"
	return nil
}
```