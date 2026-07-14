# Scenario

**Feature**: list with no channels

```
tsk channel list -> empty table or zero-count footer
```

## Steps

1. Run `tsk channel list`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "list"}
	return nil
}
```
