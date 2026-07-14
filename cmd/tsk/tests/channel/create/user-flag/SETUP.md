# Scenario

**Feature**: `create --user` sets creator participant (overrides TSK_USER)

```
TSK_USER=alice -> create --user carol -> carol + agent participants (not alice)
```

## Steps

1. Run `tsk channel create "Carol Room" --user carol`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ChannelName = "Carol Room"
	req.Args = []string{"channel", "create", req.ChannelName, "--user", "carol"}
	return nil
}
```