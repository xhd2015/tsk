# Scenario

**Feature**: `--id` is required for progress add

```
tsk progress add "text" -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"progress", "add", "some text"}
	return nil
}
```
