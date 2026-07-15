# Scenario

**Feature**: member removes another participant

```
seed -> add bob -> RemoveParticipant(bob) by alice -> alice only
```

## Steps

1. Seed; add bob; remove bob as alice.

```go
func Setup(t *testing.T, req *Request) error {
	seedChannel(t, req, "Team", "team-ch")
	store := newFileStore(t, req)
	ctx := context.Background()
	if _, err := store.AddParticipant(ctx, channel.ParticipantChangeRequest{
		ChannelID: "team-ch",
		Handle:    "bob",
		Actor:     "alice",
	}); err != nil {
		t.Fatalf("add bob: %v", err)
	}
	req.Op = "participant_remove"
	req.ChannelID = "team-ch"
	req.Sender = "alice"
	req.Handle = "bob"
	return nil
}
```