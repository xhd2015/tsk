# Scenario

**Feature**: `tsk note` help at each level

```
tsk note -h | note add -h | note list -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
