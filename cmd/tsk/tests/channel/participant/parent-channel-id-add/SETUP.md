# Scenario

**Feature**: nested `participant add` accepts parent-level `--channel-id` (peel must not break nested dispatch)

```
# nested: channel --channel-id X participant add bob
tsk channel --channel-id team-ch participant add bob
  -> peel parent opts; subcommand participant still sees "add bob"
  -> added bob\n; participants include bob
```

## Steps

1. Create channel `team-ch` (alice only).
2. Run `tsk channel --channel-id team-ch participant add bob` (no leaf `--channel-id`).

```go
func Setup(t *testing.T, req *Request) error {
	createChannel(t, req, "Team", "team-ch")
	req.Args = []string{"channel", "--channel-id", "team-ch", "participant", "add", "bob"}
	return nil
}
```
