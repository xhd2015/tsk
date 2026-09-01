# Scenario

**Feature**: `tsk advance` moves task along allowed workflow edges

```
# advance updates task.json stage only; dirname [id]-<slug> stays put
tsk advance <id> [--note N] -> stage transition create→in_process→...
```

## Preconditions

- Leaves that test advance run `create` in Setup unless testing invalid transitions.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}

// markAdvanceTree is referenced by nested intermediate SETUP packages so the
// hierarchical gen keeps a live import of this package.
func markAdvanceTree() {}
```