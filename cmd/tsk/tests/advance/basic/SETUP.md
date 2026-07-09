# Scenario

**Feature**: advance from create renames directory and updates index

```
create "add dark mode" -> advance 1 -> inbox/1-in_process-add-dark-mode/
```

## Steps

1. `tsk create "add dark mode"`.
2. `tsk advance 1`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "add dark mode"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = []string{"advance", fmt.Sprintf("%d", id)}
	return nil
}
```