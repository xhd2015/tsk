# Scenario

**Feature**: `--task` matches a title case-insensitively

```
create "Optimize Git Clone" -> search --task optimize -> 1 task match
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "Optimize Git Clone", "", nil)
	req.Args = []string{"search", "--task", "optimize"}
	return nil
}
```
