# Scenario

**Feature**: `tsk advance` moves task along allowed workflow edges

```
# advance renames <id>-<old-stage>-<slug>/ to <id>-<new-stage>-<slug>/ and updates index
tsk advance <id> [--note N] -> stage transition create→in_process→...
```

## Preconditions

- Leaves that test advance run `create` in Setup unless testing invalid transitions.

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```