# Scenario

**Feature**: create channel with slug id, creator-only participants.jsonl, index, empty messages.jsonl

```
TSK_USER=alice -> tsk channel create "Eng Alerts" -> eng-alerts
```

## Steps

1. Run `tsk channel create "Eng Alerts"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ChannelName = "Eng Alerts"
	req.Args = createChannelArgs(req.ChannelName, "")
	return nil
}
```
