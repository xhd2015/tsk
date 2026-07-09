# Scenario

**Feature**: `tsk next` returns oldest in_process task id

```
# stdout = id with earliest created_at among in_process tasks, or empty
tsk next -> <id>\n | (empty)
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```