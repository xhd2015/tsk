# Scenario

**Feature**: help documents CLI flags

```
check-channel-activity -h -> usage on stdout exit 0
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```