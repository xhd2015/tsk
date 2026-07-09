# Scenario

**Feature**: jumping create→implementation errors without mutation

```
create -> stage 1 implementation -> error; dir still *-create-*
```

## Steps

1. `tsk create "add dark mode"`.
2. `tsk stage 1 implementation` (invalid jump).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "add dark mode"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"stage", fmt.Sprintf("%d", id), "implementation"}
	return nil
}
```