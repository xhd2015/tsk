# Scenario

**Feature**: `tsk note add --id` of unknown task errors

```
tsk note add --id 99 abc -> task not found
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"note", "add", "--id", "99", "abc"}
	return nil
}
```
