# Scenario

**Feature**: `tsk channel` help surfaces subcommands and flags

```
tsk channel --help; tsk channel <subcmd> --help -> stdout usage; stderr empty
```

```go
func Setup(t *testing.T, req *Request) error {
	ensureChannelHelpersUsed()
	return nil
}
```
