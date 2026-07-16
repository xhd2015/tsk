# Scenario

**Feature**: parent and leaf `--channel-id` with different values → conflict error

```
# merge rule: set+set different -> Error: conflicting …
tsk channel --channel-id ch-a send --channel-id ch-b "msg"
  -> conflict error; no message written to ch-a
```

## Steps

1. Create channel `ch-a` (alice member). Do not create `ch-b` (conflict must fire before resolve).
2. Run send with parent `--channel-id ch-a` and leaf `--channel-id ch-b`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Channel A", "ch-a")
	req.Args = []string{
		"channel", "--channel-id", "ch-a",
		"send", "--channel-id", "ch-b", "msg",
	}
	return nil
}
```
