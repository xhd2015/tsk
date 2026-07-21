# Scenario

**Feature**: invalid stage jumps via `tsk stage` are rejected

```
# disallowed transition must not mutate task directory
tsk stage <id> <stage> -> error when edge not allowed
```

```go
func Setup(t *testing.T, req *Request) error {
	markAdvanceTree()
	ensureHelpersUsed()
	return nil
}
```
