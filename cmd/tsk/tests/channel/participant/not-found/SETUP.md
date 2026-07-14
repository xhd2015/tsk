# Scenario

**Feature**: participant ops on missing channel error

```
participant add on missing -> error
```

## Steps

1. Add participant to missing channel.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "participant", "add", "--channel-id", "missing", "bob"}
	return nil
}
```
