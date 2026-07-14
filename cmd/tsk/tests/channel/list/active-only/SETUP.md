# Scenario

**Feature**: default list hides archived channels

```
active + archived channels -> list (no --all) shows only active
```

## Steps

1. Create active and archived channels; list without `--all`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Active One", "active-one")
	createChannel(t, req, "Archived One", "archived-one")
	archiveChannel(t, req, "archived-one")
	req.Args = []string{"channel", "list"}
	return nil
}
```
