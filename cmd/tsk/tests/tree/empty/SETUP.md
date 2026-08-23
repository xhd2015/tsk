# Scenario

**Feature**: empty store prints root + `(empty)` + zero footer

```
tsk tree
```

## Steps

1. `tsk tree` on a fresh store.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"tree"}
	return nil
}
```
