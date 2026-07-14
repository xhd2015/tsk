# Scenario

**Feature**: deleted channel id cannot be recreated

```
create -> delete -> create same id -> error; tombstone remains
```

## Steps

1. Create `eng-alerts`, delete it, attempt recreate.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	deleteChannel(t, req, "eng-alerts")
	req.Args = createChannelArgs("Eng Alerts", "eng-alerts")
	return nil
}
```
