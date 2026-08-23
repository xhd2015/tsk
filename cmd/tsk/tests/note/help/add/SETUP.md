# Scenario

**Feature**: `tsk note add -h` documents `--id` and `--label`

```
tsk note add -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"note", "add", "-h"}
	return nil
}
```
