# Scenario

**Feature**: stale channel activity triggers notify command with anti-spam state

```
# last activity older than --idle (default 1h)
check-channel-activity -> idle -> exec notify script -> write state file
```

## Context

Idle leaves use `oldActivityTS` (2026-07-09) so wall-clock idle always exceeds `1h`.

```go
func Setup(t *testing.T, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```