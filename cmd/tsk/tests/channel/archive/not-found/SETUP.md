# Scenario

**Feature**: archive unknown channel id errors

```
tsk channel archive --channel-id missing -> error
```

## Steps

1. Archive nonexistent channel.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"channel", "archive", "--channel-id", "missing-ch"}
	return nil
}
```
