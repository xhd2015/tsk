# Scenario

**Feature**: archiving an already-archived channel errors

```
create -> archive -> archive again -> error
```

## Steps

1. Double archive.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Twice", "twice-ch")
	archiveChannel(t, req, "twice-ch")
	req.Args = []string{"channel", "archive", "--channel-id", "twice-ch"}
	return nil
}
```
