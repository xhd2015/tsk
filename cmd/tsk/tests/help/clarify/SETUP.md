# Scenario

**Feature**: `tsk clarify --help` lists nested subcommands

```
tsk clarify --help -> lists add, list, and confirm; exit 0
```

## Steps

1. Run `tsk clarify --help`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"clarify", "--help"}
	return nil
}
```