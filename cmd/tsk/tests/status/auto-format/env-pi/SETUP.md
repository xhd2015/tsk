# Scenario

**Feature**: `PI_CODING_AGENT` set + bare `tsk status` auto-selects agent format

```
# ExtraEnv PI_CODING_AGENT=1; no --format/--color/--plain
tsk create -> PI_CODING_AGENT=1 tsk status <id> -> agent facts
```

## Steps

1. Create inbox task.
2. Inject `PI_CODING_AGENT=1` via ExtraEnv.
3. Run bare `tsk status <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	id := createForAutoFormat(t, req, "auto env pi")
	setStatusEnv(req, "PI_CODING_AGENT=1")
	req.Args = statusArgs(id)
	return nil
}
```
