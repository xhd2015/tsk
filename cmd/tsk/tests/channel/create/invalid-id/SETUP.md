# Scenario

**Feature**: invalid channel id format rejected at create

```
tsk channel create "X" --channel-id "BAD ID" -> error
```

## Steps

1. Create with invalid `--channel-id`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = createChannelArgs("X", "BAD ID")
	return nil
}
```
