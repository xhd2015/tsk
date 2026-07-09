# Scenario

**Feature**: `tsk --help` prints top-level usage

```
tsk --help -> topHelp on stdout; exit 0
```

## Steps

1. Run `tsk --help`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--help"}
	return nil
}
```