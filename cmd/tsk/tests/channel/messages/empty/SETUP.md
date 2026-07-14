# Scenario

**Feature**: empty transcript succeeds

```
create channel -> messages with no sends -> empty or minimal stdout
```

## Steps

1. Create; messages with no sends.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Quiet", "quiet-ch")
	req.Args = []string{"channel", "messages", "--channel-id", "quiet-ch"}
	return nil
}
```
