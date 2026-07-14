# Scenario

**Feature**: forever loop checks repeatedly until signal or max ticks

```
# --forever --interval -> repeated status blocks per tick
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```