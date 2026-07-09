# Scenario

**Feature**: `tsk -h` prints top-level usage

```
tsk -h -> topHelp on stdout; exit 0
```

## Steps

1. Run `tsk -h`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	return nil
}
```