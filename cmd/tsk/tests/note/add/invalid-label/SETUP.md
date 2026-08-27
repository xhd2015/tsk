# Scenario

**Feature**: `tsk note add --label =bad` rejects empty key

```
tsk note add --id 1 --label =bad hello -> Error: empty key
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"note", "add", "--id", "1", "--label", "=bad", "hello"}
	return nil
}
```
