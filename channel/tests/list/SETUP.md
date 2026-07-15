# Scenario

**Feature**: channel Store List

```
Store.List -> active channels by default; --all includes archived
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureStoreHelpersUsed()
	return nil
}
```