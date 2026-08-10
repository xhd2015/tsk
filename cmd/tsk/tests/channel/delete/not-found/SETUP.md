# Scenario

**Feature**: delete unknown channel errors

```
tsk channel delete --channel-id nope -> error
```

## Steps

1. Delete missing channel.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"channel", "delete", "--channel-id", "nope"}
	return nil
}
```
