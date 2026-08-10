# Scenario

**Feature**: dry-run reports would notify without exec or state

```
# idle channel + --dry-run -> would notify (dry-run) -> no side effects
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```