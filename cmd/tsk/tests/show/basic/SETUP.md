# Scenario

**Feature**: show prints metadata for existing task

```
create with labels -> tsk show <id>
```

## Steps

1. `tsk create --label bug "show me"`.
2. `tsk show <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "show me"
	req.Labels = []string{"bug"}
	id := createTask(t, req, req.Title, "", req.Labels)
	req.TaskID = id
	req.Args = []string{"show", fmt.Sprintf("%d", id)}
	return nil
}
```