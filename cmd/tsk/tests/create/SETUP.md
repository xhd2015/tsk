# Scenario

**Feature**: `tsk create` allocates id, writes task directory and index entry

```
# title + optional --topic and --label flags
tsk create [--label L]... [--topic PATH] <title> -> inbox/ or topics/<path>/ task dir
```

## Preconditions

- Fresh `TSK_HOME` with no prior tasks unless a leaf Setup creates them.

## Steps

- Leaves set `req.Args` to the `create` invocation under test.

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```