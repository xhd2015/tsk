# Scenario

**Feature**: send to archived channel errors

```
archive -> send -> error
```

## Steps

1. Create, archive, attempt send.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Frozen", "frozen-ch")
	archiveChannel(t, req, "frozen-ch")
	req.Args = []string{"channel", "send", "--channel-id", "frozen-ch", "nope"}
	return nil
}
```
