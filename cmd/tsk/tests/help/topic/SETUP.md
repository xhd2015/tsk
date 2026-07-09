# Scenario

**Feature**: `tsk topic --help` lists nested subcommands

```
tsk topic --help -> lists set and mkdir; exit 0
```

## Steps

1. Run `tsk topic --help`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"topic", "--help"}
	return nil
}
```