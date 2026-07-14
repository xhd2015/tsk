# Scenario

**Feature**: participant remove without handle leaves channel

```
add bob; TSK_USER=bob; participant remove -> left team-ch\n
```

## Steps

1. Create; add bob; bob leaves.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	addParticipant(t, req, "team-ch", "bob")
	withTSKUser(t, req, "bob")
	req.Args = []string{"channel", "participant", "remove", "--channel-id", "team-ch"}
	return nil
}
```
