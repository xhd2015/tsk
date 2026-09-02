# Scenario

**Feature**: `tsk done` marks task terminal from any non-terminal stage

```
# any open stage -> done; already terminal errors; --force accepted as no-op
tsk done <id> -> stage done; further advance errors
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```
