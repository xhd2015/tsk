# Scenario

**Feature**: `tsk topic alias add` writes topic.json

```
tsk topic alias add <topic> <alias>
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
