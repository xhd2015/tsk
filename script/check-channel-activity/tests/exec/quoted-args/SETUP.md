# Scenario

**Feature**: LINE with quoted argument containing spaces

```
# LINE: /bin/sh -c 'echo "$1" > argv.txt' "hello world"
# idle -> exec -> argv.txt contains hello world
```

## Steps

1. Seed idle channel.
2. Run with LINE embedding quoted `"hello world"` argument.

```go
func Setup(t *testing.T, req *Request) error {
	msgs := []channelMessage{{
		ID: 1, Sender: "alice", Body: "stale", CreatedAt: oldActivityTS,
	}}
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, msgs)
	req.ArgvPath = filepath.Join(req.WorkRoot, "argv.txt")
	execLine := fmt.Sprintf(`/bin/sh -c 'echo "$1" > %s' "hello world"`, req.ArgvPath)
	req.Args = []string{
		"--channel-id", req.ChannelID,
		"--exec-if-idle-1h", execLine,
	}
	return nil
}
```