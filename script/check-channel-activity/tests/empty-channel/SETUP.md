# Scenario

**Feature**: empty messages.jsonl uses channel.created_at as last activity

```
# no messages; old created_at -> idle from created_at -> notify
```

## Steps

1. Seed active channel with empty messages and old `created_at`.
2. Run one-shot check.

```go
func Setup(t *testing.T, req *Request) error {
	req.LastActivity = writeActiveChannel(t, req, oldCreatedAtTS, nil)
	req.Args = defaultCheckArgs(req)
	return nil
}
```