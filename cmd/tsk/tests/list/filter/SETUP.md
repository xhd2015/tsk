# Scenario

**Feature**: list --stage create returns only create-stage tasks

```
create A; create B + advance B -> list --stage create -> stdout "1\n"
```

## Steps

1. Create task "stay".
2. Create task "go" and advance to `in_process`.
3. `tsk list --stage create`.

```go
func Setup(t *testing.T, req *Request) error {
	id1 := createTask(t, req, "stay", "", nil)
	id2 := createTask(t, req, "go", "", nil)
	advanceTask(t, req, id2, "")
	req.TaskID = id1
	req.Args = []string{"list", "--stage", "create"}
	return nil
}
```