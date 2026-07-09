# Scenario

**Feature**: `--color` forces diagram path (blocks auto-agent) even under CODEX

```
# ExtraEnv CODEX_THREAD_ID=t1; argv --color
tsk create -> CODEX_THREAD_ID=t1 tsk status --color <id> -> diagram (may ANSI), not agent facts
```

## Steps

1. Create inbox task.
2. Inject `CODEX_THREAD_ID=t1`.
3. Run `tsk status --color <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto force color")
	setStatusEnv(req, "CODEX_THREAD_ID=t1")
	req.Args = statusArgs(id, "--color")
	return nil
}
```
