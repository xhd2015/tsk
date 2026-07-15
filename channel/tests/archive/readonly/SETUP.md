# Scenario

**Feature**: send blocked on archived channel

```
seed -> archive -> send -> error; no messages
```

## Steps

1. Seed channel; archive; attempt send.

```go
func Setup(t *testing.T, req *Request) error {
	seedChannel(t, req, "Readonly", "readonly-ch")
	req.Op = "archive_readonly"
	req.ChannelID = "readonly-ch"
	req.Sender = "alice"
	req.MessageBody = "nope"
	return nil
}
```