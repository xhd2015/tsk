# Scenario

**Feature**: `tsk status --help` documents the `--format` flag

```
tsk status --help -> usage includes --format; exit 0
```

## Steps

1. Run `tsk status --help`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"status", "--help"}
	return nil
}
```
