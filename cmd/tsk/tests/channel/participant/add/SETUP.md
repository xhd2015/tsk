# Scenario

**Feature**: add participant prints added handle

```
participant add bob -> added bob\n; bob in channel.json
```

## Steps

1. Create; add bob.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	req.Args = []string{"channel", "participant", "add", "--channel-id", "team-ch", "bob"}
	return nil
}
```
