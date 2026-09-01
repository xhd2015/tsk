# Scenario

**Feature**: `tsk add --note` appends task notes in the same invocation

```
# significance: one-shot add + note(s); topic/inbox/parent are placement variants
tsk add [--topic|--parent] [--note T]... <title>
```

## Preconditions

- Fresh `TSK_HOME` from root setup.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
