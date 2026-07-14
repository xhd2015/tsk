# Scenario

**Feature**: `list --json` emits JSON array without ANSI

```
create channel -> list --json -> valid JSON array
```

## Steps

1. Create one channel; `tsk channel list --json`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "JSON Room", "json-room")
	req.Args = []string{"channel", "list", "--json"}
	return nil
}
```
