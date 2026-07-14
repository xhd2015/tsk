# Scenario

**Feature**: `list --all` includes archived channels

```
active + archived -> list --all shows both
```

## Steps

1. Seed channels; run `tsk channel list --all`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Active One", "active-one")
	createChannel(t, req, "Archived One", "archived-one")
	archiveChannel(t, req, "archived-one")
	req.Args = []string{"channel", "list", "--all"}
	return nil
}
```
