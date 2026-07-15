# Scenario

**Feature**: cannot remove the last remaining participant

```
creator alice only -> attempt remove alice -> error
```

## Steps

1. Create channel with sole creator; attempt remove alice.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Solo", "solo-ch")
	req.Args = []string{"channel", "participant", "remove", "--channel-id", "solo-ch", "alice"}
	return nil
}
```