# Scenario

**Feature**: `TSK_STATUS_FORMAT=agent` forces agent format with host env cleared

```
# ExtraEnv TSK_STATUS_FORMAT=agent only (no CODEX/PI)
tsk create -> TSK_STATUS_FORMAT=agent tsk status <id> -> agent facts
```

## Steps

1. Create inbox task.
2. Inject only `TSK_STATUS_FORMAT=agent` (host detect vars remain cleared by tskEnv).
3. Run bare `tsk status <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto tsk format agent")
	setStatusEnv(req, "TSK_STATUS_FORMAT=agent")
	req.Args = statusArgs(id)
	return nil
}
```
