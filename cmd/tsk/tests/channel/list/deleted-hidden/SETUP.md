# Scenario

**Feature**: deleted (tombstoned) channels never appear in list --all

```
create -> delete -> list --all -> channel absent
```

## Steps

1. Create and delete channel; list with `--all`.

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Gone", "gone-channel")
	deleteChannel(t, req, "gone-channel")
	req.Args = []string{"channel", "list", "--all"}
	return nil
}
```
