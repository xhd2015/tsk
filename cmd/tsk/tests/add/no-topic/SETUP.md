# Scenario

**Feature**: create task without --topic lands in inbox

```
tsk add "add dark mode" -> inbox/[1]-create-add-dark-mode/
```

## Steps

1. Run `tsk add "add dark mode"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "add dark mode"
	req.Args = []string{"add", req.Title}
	return nil
}
```