# Scenario

**Feature**: `tsk create --help` prints create usage and flags

```
tsk create --help -> create usage with --label and --topic; exit 0
```

## Steps

1. Run `tsk create --help`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"create", "--help"}
	return nil
}
```