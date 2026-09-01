# Scenario

**Feature**: successful create appends audit event

```
tsk add "audit me" -> events.jsonl line for add
```

## Steps

1. Run `tsk add "audit me"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "audit me"
	req.Args = []string{"add", req.Title}
	return nil
}
```