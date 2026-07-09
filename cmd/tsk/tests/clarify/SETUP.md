# Scenario

**Feature**: `tsk clarify` manages clarification batch during clarification stage

```
# add questions to clarify/batch.json; confirm -y advances to implementation
tsk clarify add|list|confirm <id> -> batch.json + stage rename on confirm -y
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```