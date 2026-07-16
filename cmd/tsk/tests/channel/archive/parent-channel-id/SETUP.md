# Scenario

**Feature**: parent-level `--channel-id` before `archive` archives the channel

```
# parent peel for lifecycle
tsk channel --channel-id eng-alerts archive
  -> archived eng-alerts\n; active -> archive/
```

## Steps

1. Create channel `eng-alerts`.
2. Run `tsk channel --channel-id eng-alerts archive` (no leaf `--channel-id`).

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	req.Args = []string{"channel", "--channel-id", "eng-alerts", "archive"}
	return nil
}
```
