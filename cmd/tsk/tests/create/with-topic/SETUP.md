# Scenario

**Feature**: create task with --topic places directory under topics tree

```
tsk create --topic eng/backend "x" -> topics/eng/backend/1-create-x/
```

## Steps

1. Run `tsk create --topic eng/backend "x"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "x"
	req.Topic = "eng/backend"
	req.Args = []string{"create", "--topic", req.Topic, req.Title}
	return nil
}
```