# Scenario

**Feature**: channel commands append `events.jsonl` audit lines

```
tsk channel <subcmd> -> events.jsonl line with command: channel
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
