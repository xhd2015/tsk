# Scenario

**Feature**: channel Store Create

```
Creator -> Store.Create -> active/<id>/ with normalized layout
```

## Steps

1. Configure `req.Op=create` and channel name/id in leaf Setup.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureStoreHelpersUsed()
	return nil
}
```