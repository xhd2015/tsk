# Scenario

**Feature**: `tsk show` prints task metadata block

```
tsk show <id> -> title, stage, labels, topic, timestamps
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```