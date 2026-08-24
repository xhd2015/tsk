# Scenario

**Feature**: `tsk delete -h` documents usage and `--recursive`

```
tsk delete -h -> Usage: tsk delete [--recursive] <id>
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"delete", "-h"}
	return nil
}
```
