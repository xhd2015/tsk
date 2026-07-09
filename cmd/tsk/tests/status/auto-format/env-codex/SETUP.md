# Scenario

**Feature**: `CODEX_THREAD_ID` set + bare `tsk status` auto-selects agent format

```
# ExtraEnv CODEX_THREAD_ID=t1; no --format/--color/--plain
tsk create -> CODEX_THREAD_ID=t1 tsk status <id> -> agent facts (id/title/topic/dir)
```

## Steps

1. Create inbox task.
2. Inject `CODEX_THREAD_ID=t1` via ExtraEnv.
3. Run bare `tsk status <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto env codex")
	setStatusEnv(req, "CODEX_THREAD_ID=t1")
	req.Args = statusArgs(id)
	return nil
}
```
