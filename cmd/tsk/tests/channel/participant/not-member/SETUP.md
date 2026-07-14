# Scenario

**Feature**: non-participant cannot add others

```
alice channel; charlie tries participant add -> error
```

## Steps

1. Create; charlie tries add.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Closed", "closed-ch")
	withTSKUser(t, req, "charlie")
	req.Args = []string{"channel", "participant", "add", "--channel-id", "closed-ch", "dave"}
	return nil
}
```
