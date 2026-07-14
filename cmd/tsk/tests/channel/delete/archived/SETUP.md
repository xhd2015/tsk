# Scenario

**Feature**: delete works on archived channels

```
create -> archive -> delete -> tombstone
```

## Steps

1. Create, archive, delete.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Old", "old-ch")
	archiveChannel(t, req, "old-ch")
	req.Args = []string{"channel", "delete", "--channel-id", "old-ch"}
	return nil
}
```
