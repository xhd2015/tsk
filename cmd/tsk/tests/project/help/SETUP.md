# Scenario

**Feature**: project help surfaces at command and leaf levels

```
tsk project -h; tsk project add -h; tsk project tree -h; tsk project list -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureProjectHelpersUsed()
	req.Args = []string{"project", "-h"}
	return nil
}
```
