# Scenario

**Feature**: duplicate channel id rejected

```
Create eng-alerts -> Create eng-alerts again -> error
```

## Steps

1. Run create twice with same id.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "create_duplicate"
	req.ChannelName = "Eng Alerts"
	req.ChannelID = "eng-alerts"
	return nil
}
```