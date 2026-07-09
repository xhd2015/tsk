# Scenario

**Feature**: CLI UX — single stderr error path and create prints task id

```
# errors: one line on stderr via single path; create success prints id + newline on stdout
tsk advance -> stderr once; tsk create "title" -> stdout id\n + inbox side effects
```

## Preconditions

- `fail()` must not duplicate errors already printed by `main`.
- `create` success stdout is task id and trailing newline only.

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```