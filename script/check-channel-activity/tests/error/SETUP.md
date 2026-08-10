# Scenario

**Feature**: checker errors on missing or archived channels

```
# missing index/active dir OR archive/ path -> Error: ... exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```