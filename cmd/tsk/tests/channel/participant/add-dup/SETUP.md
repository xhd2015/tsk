# Scenario

**Feature**: adding existing participant is idempotent

```
add bob twice -> both succeed; roster unchanged count
```

## Steps

1. Create; add bob; add bob again.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	addParticipant(t, req, "team-ch", "bob")
	req.Args = []string{"channel", "participant", "add", "--channel-id", "team-ch", "bob"}
	return nil
}
```
