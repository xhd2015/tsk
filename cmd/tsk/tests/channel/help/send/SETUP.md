# Scenario

**Feature**: `channel send --help` documents --channel-id and --user

```
tsk channel send --help -> --channel-id, --user
```

## Steps

1. Run send help.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "send", "--help"}
	return nil
}
```