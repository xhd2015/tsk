# Scenario

**Feature**: archived channels are readonly for mutations

```
Archive -> SendMessage -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureStoreHelpersUsed()
	return nil
}
```