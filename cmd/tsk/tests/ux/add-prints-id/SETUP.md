# Scenario

**Feature**: `tsk add` prints allocated task id on stdout

```
tsk add "hello" -> stdout "1\n"; inbox task dir created
```

## Steps

1. Run `tsk add "hello"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "hello"
	req.Args = []string{"add", req.Title}
	return nil
}
```