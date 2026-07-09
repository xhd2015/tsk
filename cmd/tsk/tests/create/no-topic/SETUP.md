# Scenario

**Feature**: create task without --topic lands in inbox

```
tsk create "add dark mode" -> inbox/1-create-add-dark-mode/
```

## Steps

1. Run `tsk create "add dark mode"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "add dark mode"
	req.Args = []string{"create", req.Title}
	return nil
}
```