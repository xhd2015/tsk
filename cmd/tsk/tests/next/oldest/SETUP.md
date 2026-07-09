# Scenario

**Feature**: next picks older of two in_process tasks

```
create A -> create B -> advance both -> tsk next -> stdout "1"
```

## Steps

1. Create task "older".
2. Create task "newer".
3. Advance both to `in_process`.
4. Run `tsk next`.

```go
func Setup(t *testing.T, req *Request) error {
	id1 := createTask(t, req, "older", "", nil)
	id2 := createTask(t, req, "newer", "", nil)
	advanceTask(t, req, id1, "")
	advanceTask(t, req, id2, "")
	req.TaskID = id1
	req.Args = []string{"next"}
	return nil
}
```