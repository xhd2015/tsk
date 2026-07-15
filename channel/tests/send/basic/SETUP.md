# Scenario

**Feature**: participant sends first message

```
seed channel -> SendMessage(alice, "hello") -> id 1 in jsonl
```

## Steps

1. Seed channel `send-ch`.
2. Send message as alice.

```go
func Setup(t *testing.T, req *Request) error {
	seedChannel(t, req, "Send Room", "send-ch")
	req.Op = "send"
	req.ChannelID = "send-ch"
	req.Sender = "alice"
	req.MessageBody = "hello team"
	return nil
}
```