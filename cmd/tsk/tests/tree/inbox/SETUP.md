# Scenario

**Feature**: inbox-only tasks appear at root level

```
create inbox task -> tree
```

## Steps

1. `tsk add "solo"`.
2. `tsk tree`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "solo", "", nil)
	req.Args = []string{"tree"}
	return nil
}
```
