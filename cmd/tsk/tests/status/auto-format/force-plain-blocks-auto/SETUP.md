# Scenario

**Feature**: `--plain` forces diagram path (blocks auto-agent) even under CODEX

```
# ExtraEnv CODEX_THREAD_ID=t1; argv --plain
tsk create -> CODEX_THREAD_ID=t1 tsk status --plain <id> -> ASCII diagram, not agent facts
```

## Steps

1. Create inbox task.
2. Inject `CODEX_THREAD_ID=t1`.
3. Run `tsk status --plain <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto force plain")
	setStatusEnv(req, "CODEX_THREAD_ID=t1")
	req.Args = statusArgs(id, "--plain")
	return nil
}
```
