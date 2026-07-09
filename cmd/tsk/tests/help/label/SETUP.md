# Scenario

**Feature**: `tsk label --help` lists nested subcommands

```
tsk label --help -> lists add and rm; exit 0
```

## Steps

1. Run `tsk label --help`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"label", "--help"}
	return nil
}
```