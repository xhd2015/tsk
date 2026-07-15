# Scenario

**Feature**: archived channels are readonly for mutations

```
Archive -> SendMessage -> error
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureStoreHelpersUsed()
	return nil
}
```