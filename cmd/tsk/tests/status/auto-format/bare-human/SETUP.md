# Scenario

**Feature**: bare `tsk status` with no agent host env selects diagram (not agent facts)

```
# ExtraEnv empty; tskEnv cleared CODEX/PI/TSK_STATUS_FORMAT
tsk create "auto bare human" -> tsk status <id> -> diagram box art
```

## Steps

1. Create inbox task (stage `create`).
2. Run bare `tsk status <id>` with no format flags and no host-agent ExtraEnv.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto bare human")
	// ExtraEnv intentionally empty: host agent vars already stripped by tskEnv
	req.Args = statusArgs(id)
	return nil
}
```
