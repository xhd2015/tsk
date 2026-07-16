# Scenario

**Feature**: parent-level `--channel-id` before `messages` lists the channel transcript

```
# parent peel for read path
tsk channel --channel-id chat-ch messages
  -> same transcript as leaf --channel-id form
```

## Steps

1. Create channel; send two messages via leaf form.
2. Run `tsk channel --channel-id chat-ch messages`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Chat", "chat-ch")
	sendChannelMessage(t, req, "chat-ch", "first")
	sendChannelMessage(t, req, "chat-ch", "second")
	req.Args = []string{"channel", "--channel-id", "chat-ch", "messages"}
	return nil
}
```
