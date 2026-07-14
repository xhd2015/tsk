# Scenario

**Feature**: remove other participant by handle

```
add bob -> participant remove bob -> removed bob\n
```

## Steps

1. Create; add bob; alice removes bob.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	addParticipant(t, req, "team-ch", "bob")
	req.Args = []string{"channel", "participant", "remove", "--channel-id", "team-ch", "bob"}
	return nil
}
```
