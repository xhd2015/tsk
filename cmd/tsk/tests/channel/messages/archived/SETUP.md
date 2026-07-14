# Scenario

**Feature**: participants can read messages on archived channels

```
create -> send -> archive -> messages still works
```

## Steps

1. Create, send, archive, read messages.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "History", "hist-ch")
	sendChannelMessage(t, req, "hist-ch", "preserved")
	archiveChannel(t, req, "hist-ch")
	req.Args = []string{"channel", "messages", "--channel-id", "hist-ch"}
	return nil
}
```
