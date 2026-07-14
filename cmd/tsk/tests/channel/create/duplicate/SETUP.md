# Scenario

**Feature**: second create with same channel id errors

```
create eng-alerts -> create eng-alerts again -> error
```

## Steps

1. Create channel `eng-alerts`.
2. Attempt duplicate create with same `--channel-id`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	req.Args = createChannelArgs("Eng Alerts Again", "eng-alerts")
	return nil
}
```
