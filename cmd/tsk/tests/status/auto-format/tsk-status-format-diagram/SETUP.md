# Scenario

**Feature**: `TSK_STATUS_FORMAT=diagram` overrides CODEX host detect → diagram

```
# ExtraEnv CODEX_THREAD_ID=t1 + TSK_STATUS_FORMAT=diagram
# TSK_STATUS_FORMAT beats detect; no --format flag
tsk create -> env tsk status <id> -> diagram (not agent)
```

## Steps

1. Create inbox task.
2. Inject `CODEX_THREAD_ID=t1` and `TSK_STATUS_FORMAT=diagram`.
3. Run bare `tsk status <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto tsk format diagram")
	setStatusEnv(req,
		"CODEX_THREAD_ID=t1",
		"TSK_STATUS_FORMAT=diagram",
	)
	req.Args = statusArgs(id)
	return nil
}
```
