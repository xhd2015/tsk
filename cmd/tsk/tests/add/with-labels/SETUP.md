# Scenario

**Feature**: create with multiple --label flags stores sorted labels

```
tsk add --label bug --label urgent "x" -> task.json labels ["bug","urgent"]
```

## Steps

1. Run `tsk add --label bug --label urgent "x"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "x"
	req.Labels = []string{"bug", "urgent"}
	req.Args = addArgs(req.Title, "", req.Labels)
	return nil
}
```