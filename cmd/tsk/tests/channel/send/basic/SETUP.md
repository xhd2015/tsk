# Scenario

**Feature**: participant sends message appended to messages.jsonl

```
alice member -> send "fix login" -> sent message 1\n; jsonl line
```

## Steps

1. Create channel; send message.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	req.Args = []string{"channel", "send", "--channel-id", "eng-alerts", "fix login"}
	return nil
}
```
