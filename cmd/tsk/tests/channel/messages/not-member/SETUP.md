# Scenario

**Feature**: non-participant cannot read messages

```
alice creates; bob reads messages -> error
```

## Steps

1. Create; switch to bob; messages.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Secret", "secret-ch")
	withTSKUser(t, req, "bob")
	req.Args = []string{"channel", "messages", "--channel-id", "secret-ch"}
	return nil
}
```
