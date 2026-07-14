# Scenario

**Feature**: archived channels reject send and participant mutations

```
archive channel -> send / participant add / participant remove -> errors
```

## Steps

1. Create, archive; attempt `send`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Readonly", "readonly-ch")
	archiveChannel(t, req, "readonly-ch")
	req.Args = []string{"channel", "send", "--channel-id", "readonly-ch", "hello"}
	return nil
}
```
