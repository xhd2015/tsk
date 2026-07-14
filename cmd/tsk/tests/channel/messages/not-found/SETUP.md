# Scenario

**Feature**: messages on missing channel errors

```
messages --channel-id missing -> error
```

## Steps

1. Messages without channel.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "messages", "--channel-id", "missing"}
	return nil
}
```
