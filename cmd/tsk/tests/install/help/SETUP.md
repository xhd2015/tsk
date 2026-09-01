# Scenario

**Feature**: `tsk install --help` lists available wrappers

```
tsk install --help -> usage + pmark; exit 0
```

## Steps

1. Run `tsk install --help`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"install", "--help"}
	return nil
}
```
