# Scenario

**Feature**: topic nested help

```
tsk topic where --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
