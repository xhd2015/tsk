# Scenario

**Feature**: `tsk` with no arguments prints top-level usage

```
tsk -> topHelp on stdout; exit 0
```

## Steps

1. Run `tsk` with no arguments.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = nil
	return nil
}
```