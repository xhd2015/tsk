# Scenario

**Feature**: create with explicit --channel-id

```
tsk channel create "Room" --channel-id my-room -> my-room
```

## Steps

1. Run create with `--channel-id my-room`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ChannelName = "Room"
	req.ChannelID = "my-room"
	req.Args = createChannelArgs(req.ChannelName, req.ChannelID)
	return nil
}
```
