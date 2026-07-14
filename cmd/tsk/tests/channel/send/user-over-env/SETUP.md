# Scenario

**Feature**: `--user` flag wins over `TSK_USER` env for send identity

```
TSK_USER=alice -> send --user bob -> sender bob (not alice)
```

## Steps

1. Create as alice (default TSK_USER); add bob; send with `--user bob`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Override", "override-ch")
	addParticipant(t, req, "override-ch", "bob")
	req.Args = []string{"channel", "send", "--channel-id", "override-ch", "--user", "bob", "flag wins"}
	return nil
}
```