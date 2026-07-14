# Scenario

**Feature**: non-participant cannot send

```
alice creates; bob (TSK_USER) tries send -> error
```

## Steps

1. Create as alice; switch to bob; send.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Private", "private-ch")
	withTSKUser(t, req, "bob")
	req.Args = []string{"channel", "send", "--channel-id", "private-ch", "hi"}
	return nil
}
```
