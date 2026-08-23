# Scenario

**Feature**: `tsk create` allocates id, writes task directory and index entry

```
# title + optional --topic / --parent / --label / --note flags
tsk create [--label L]... [--topic PATH | --parent ID] [--note T]... <title>
  -> inbox/, topics/<path>/, or nested under parent; optional notes.jsonl entries
```

## Preconditions

- Fresh `TSK_HOME` with no prior tasks unless a leaf Setup creates them.

## Steps

- Leaves set `req.Args` to the `create` invocation under test.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```