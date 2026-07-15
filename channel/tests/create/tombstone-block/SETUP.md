# Scenario

**Feature**: tombstone blocks id reuse

```
Create -> Delete -> Create same id -> error; tombstone remains
```

## Steps

1. Create, delete, recreate same id.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "create_tombstone_block"
	req.ChannelName = "Eng Alerts"
	req.ChannelID = "eng-alerts"
	return nil
}
```