# Scenario

**Feature**: forever loop checks repeatedly until signal or max ticks

```
# --forever --interval -> repeated status blocks per tick
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```