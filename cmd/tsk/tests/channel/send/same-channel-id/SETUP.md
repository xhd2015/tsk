# Scenario

**Feature**: parent and leaf `--channel-id` with the same value is OK (no conflict)

```
# merge rule: set+set same -> OK
tsk channel --channel-id eng-alerts send --channel-id eng-alerts "hi"
  -> use eng-alerts -> sent message 1
```

## Steps

1. Create channel `eng-alerts`.
2. Run send with parent and leaf `--channel-id` both `eng-alerts`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	req.Args = []string{
		"channel", "--channel-id", "eng-alerts",
		"send", "--channel-id", "eng-alerts", "hi",
	}
	return nil
}
```
