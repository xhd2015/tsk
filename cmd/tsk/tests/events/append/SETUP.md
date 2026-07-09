# Scenario

**Feature**: successful create appends audit event

```
tsk create "audit me" -> events.jsonl line for create
```

## Steps

1. Run `tsk create "audit me"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "audit me"
	req.Args = []string{"create", req.Title}
	return nil
}
```