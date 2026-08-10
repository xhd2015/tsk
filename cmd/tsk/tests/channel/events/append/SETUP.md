# Scenario

**Feature**: channel create appends events.jsonl audit line

```
tsk channel create -> events.jsonl command channel
```

## Steps

1. Run channel create.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = createChannelArgs("Audit Channel", "audit-ch")
	return nil
}
```
