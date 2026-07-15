# Scenario

**Feature**: `participants --json` lists roster for members

```
create -> participants --json -> JSON array with alice only
```

## Steps

1. Create; participants --json.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Roster", "roster-ch")
	req.Args = []string{"channel", "participants", "--channel-id", "roster-ch", "--json"}
	return nil
}
```