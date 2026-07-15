# Scenario

**Feature**: invalid channel id format rejected

```
Create with id "BAD ID" -> error
```

## Steps

1. Create with invalid id.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "create_invalid_id"
	req.ChannelName = "Bad"
	req.ChannelID = "BAD ID"
	return nil
}
```