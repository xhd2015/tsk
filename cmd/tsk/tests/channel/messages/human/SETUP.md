# Scenario

**Feature**: human transcript in chronological order

```
send two messages -> messages shows [id] sender timestamp + body blocks
```

## Steps

1. Create; send two messages; run `messages`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Chat", "chat-ch")
	sendChannelMessage(t, req, "chat-ch", "first")
	sendChannelMessage(t, req, "chat-ch", "second")
	req.Args = []string{"channel", "messages", "--channel-id", "chat-ch"}
	return nil
}
```
