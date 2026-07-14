# Scenario

**Feature**: participant add/remove blocked on archived channel

```
archive -> participant add -> error
```

## Steps

1. Create, archive; attempt add.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Frozen", "frozen-ch")
	archiveChannel(t, req, "frozen-ch")
	req.Args = []string{"channel", "participant", "add", "--channel-id", "frozen-ch", "bob"}
	return nil
}
```
