# Scenario

**Bug**: duplicate error lines when `tsk advance` is missing task id

```
tsk advance -> stderr contains "task id required" exactly once; exit 1
```

## Steps

1. Run `tsk advance` with no task id.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"advance"}
	return nil
}
```