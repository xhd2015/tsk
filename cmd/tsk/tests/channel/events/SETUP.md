# Scenario

**Feature**: channel commands append `events.jsonl` audit lines

```
tsk channel <subcmd> -> events.jsonl line with command: channel
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
