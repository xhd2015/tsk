# Scenario

**Feature**: create with multiple --label flags stores sorted labels

```
tsk create --label bug --label urgent "x" -> task.json labels ["bug","urgent"]
```

## Steps

1. Run `tsk create --label bug --label urgent "x"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "x"
	req.Labels = []string{"bug", "urgent"}
	req.Args = createArgs(req.Title, "", req.Labels)
	return nil
}
```