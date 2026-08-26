# Scenario

**Feature**: `--color` and `--no-color` together error

```
tsk search --color --no-color foo -> conflict error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"search", "--color", "--no-color", "foo"}
	return nil
}
```
