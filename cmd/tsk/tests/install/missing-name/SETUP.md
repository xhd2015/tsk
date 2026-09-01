# Scenario

**Feature**: `tsk install` without a name errors

```
tsk install -> Error: name required; non-zero
```

## Steps

1. Run `tsk install` with no name.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"install"}
	return nil
}
```
