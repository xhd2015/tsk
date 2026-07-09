# Scenario

**Feature**: `--format=diagram` wins over CODEX host detect

```
# ExtraEnv CODEX_THREAD_ID=t1; argv --format=diagram
tsk create -> CODEX_THREAD_ID=t1 tsk status --format=diagram <id> -> diagram
```

## Steps

1. Create inbox task.
2. Inject `CODEX_THREAD_ID=t1`.
3. Run `tsk status --format=diagram <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto force diagram flag")
	setStatusEnv(req, "CODEX_THREAD_ID=t1")
	req.Args = statusArgs(id, "--format=diagram")
	return nil
}
```
