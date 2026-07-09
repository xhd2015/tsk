# Scenario

**Feature**: `tsk done` marks task terminal from allowed stages

```
# summary -> done (not via advance); user_followup -> done also allowed
tsk done <id> -> stage done; further advance errors
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```