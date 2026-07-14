# Scenario

**Feature**: `messages --json` returns JSON array

```
send message -> messages --json -> array with message object
```

## Steps

1. Create, send, messages --json.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Data", "data-ch")
	sendChannelMessage(t, req, "data-ch", "json body")
	req.Args = []string{"channel", "messages", "--channel-id", "data-ch", "--json"}
	return nil
}
```
