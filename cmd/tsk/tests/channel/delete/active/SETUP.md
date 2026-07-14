# Scenario

**Feature**: delete active channel writes tombstone and removes from list

```
create -> delete -> deleted id\n; tombstone; not in list --all
```

## Steps

1. Create and delete active channel.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Temp", "temp-ch")
	req.Args = []string{"channel", "delete", "--channel-id", "temp-ch"}
	return nil
}
```
