# Scenario

**Feature**: `tsk followup` records user message and sets user_followup stage

```
# from summary only; writes context/followup-<ts>.md
tsk followup <id> <message> -> user_followup stage
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```