# Scenario

**Feature**: create task with --topic places directory under topics tree

```
tsk add --topic eng/backend "x" -> topics/eng/backend/[1]-create-x/
```

## Steps

1. Run `tsk add --topic eng/backend "x"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "x"
	req.Topic = "eng/backend"
	req.Args = []string{"add", "--topic", req.Topic, req.Title}
	return nil
}
```