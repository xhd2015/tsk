# Scenario

**Feature**: `--limit 1` returns only the last message

```
send two -> messages --limit 1 -> only message 2
```

## Steps

1. Create; send two; messages --limit 1.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Limit", "limit-ch")
	sendChannelMessage(t, req, "limit-ch", "one")
	sendChannelMessage(t, req, "limit-ch", "two")
	req.Args = []string{"channel", "messages", "--channel-id", "limit-ch", "--limit", "1"}
	return nil
}
```
