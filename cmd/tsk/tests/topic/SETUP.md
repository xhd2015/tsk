# Scenario

**Feature**: `tsk topic` manages topic tree and task placement

```
# topic mkdir creates path; topic set moves task dir and updates topic_path + index
tsk topic mkdir <path>
tsk topic set <id> <path|--inbox>
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```