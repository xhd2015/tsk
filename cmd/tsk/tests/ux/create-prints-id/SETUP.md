# Scenario

**Feature**: `tsk create` prints allocated task id on stdout

```
tsk create "hello" -> stdout "1\n"; inbox task dir created
```

## Steps

1. Run `tsk create "hello"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "hello"
	req.Args = []string{"create", req.Title}
	return nil
}
```