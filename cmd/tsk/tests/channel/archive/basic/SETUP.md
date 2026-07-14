# Scenario

**Feature**: archive moves active dir to archive and updates index

```
create eng-alerts -> archive -> archived eng-alerts\n; excluded from default list
```

## Steps

1. Create channel; run `tsk channel archive --channel-id eng-alerts`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	req.Args = []string{"channel", "archive", "--channel-id", "eng-alerts"}
	return nil
}
```
