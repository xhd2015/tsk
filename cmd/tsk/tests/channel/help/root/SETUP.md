# Scenario

**Feature**: `tsk channel --help` lists subcommands

```
tsk channel --help -> create, list, archive, delete, send, messages, participants, participant
```

## Steps

1. Run channel root help.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"channel", "--help"}
	return nil
}
```
