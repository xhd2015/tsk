# Scenario

**Feature**: send to missing channel errors

```
send --channel-id missing -> error
```

## Steps

1. Send without creating channel.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "send", "--channel-id", "missing", "hi"}
	return nil
}
```
