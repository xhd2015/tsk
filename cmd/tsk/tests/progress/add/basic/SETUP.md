# Scenario

**Feature**: add a progress entry with required status

```
create task -> progress add --status in-progress -> added progress
```

## Steps

1. `tsk add "demo"`.
2. `tsk progress add --id 1 --status in-progress "investigating"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	req.Args = []string{"progress", "add", "--id", "1", "--status", "in-progress", "investigating"}
	return nil
}
```
