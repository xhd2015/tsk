# Scenario

**Feature**: SIGINT/SIGTERM stops forever loop gracefully

```
# --forever + signal -> stderr "stopped\n" -> exit 0
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```