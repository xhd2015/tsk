# Scenario

**Feature**: `list` hard-rejects parent-level `--channel-id`

```
# list does not accept shared parent flags
tsk channel --channel-id eng-alerts list
  -> Error: … --channel-id not accepted … (or equivalent)
```

## Steps

1. Create a channel so list would otherwise succeed without the flag.
2. Run `tsk channel --channel-id eng-alerts list`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Eng Alerts", "eng-alerts")
	req.Args = []string{"channel", "--channel-id", "eng-alerts", "list"}
	return nil
}
```
