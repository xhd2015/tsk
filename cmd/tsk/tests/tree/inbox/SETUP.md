# Scenario

**Feature**: inbox-only tasks appear at root level

```
create inbox task -> tree
```

## Steps

1. `tsk create "solo"`.
2. `tsk tree`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "solo", "", nil)
	req.Args = []string{"tree"}
	return nil
}
```
