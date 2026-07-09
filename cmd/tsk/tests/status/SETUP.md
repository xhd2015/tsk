# Scenario

**Feature**: `tsk status` prints ASCII pipeline diagram for a task

```
# NOT a list filter; shows workflow stages with marker on current stage
tsk status <id> -> pipeline diagram stdout
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```