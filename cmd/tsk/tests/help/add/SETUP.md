# Scenario

**Feature**: `tsk add --help` prints add usage and flags

```
tsk add --help -> add usage with --label and --topic; exit 0
```

## Steps

1. Run `tsk add --help`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"add", "--help"}
	return nil
}
```