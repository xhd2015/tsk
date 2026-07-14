# Scenario

**Feature**: `channel create --help` documents --channel-id

```
tsk channel create --help -> --channel-id flag
```

## Steps

1. Run create help.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "create", "--help"}
	return nil
}
```
