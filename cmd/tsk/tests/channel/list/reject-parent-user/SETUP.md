# Scenario

**Feature**: `list` hard-rejects parent-level `--user`

```
# list does not accept shared parent flags
tsk channel --user alice list
  -> Error: … --user not accepted … (or equivalent)
```

## Steps

1. Run `tsk channel --user alice list` (empty or seeded home both OK; reject is parse-time).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "--user", "alice", "list"}
	return nil
}
```
