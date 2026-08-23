# Scenario

**Feature**: `tsk topic notes --help`

```
tsk topic notes --help
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"topic", "notes", "--help"}
	return nil
}
```
