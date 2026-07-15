# Scenario

**Feature**: create channel with creator-only participants.jsonl

```
Creator=alice -> Create("Eng Alerts") -> eng-alerts active dir
```

## Steps

1. Create channel `Eng Alerts` with default slug id.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "create"
	req.ChannelName = "Eng Alerts"
	req.ChannelID = ""
	return nil
}
```