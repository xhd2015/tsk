# Scenario

**Feature**: `tsk list` prints task ids with optional filters

```
tsk list [--stage S] [--label L] [--topic PREFIX] -> one id per line
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```