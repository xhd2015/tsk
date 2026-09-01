# Scenario

**Feature**: `tsk install` rejects unknown names

```
tsk install foo -> Error: unknown name; non-zero
```

## Steps

1. Run `tsk install foo`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"install", "foo"}
	return nil
}
```
