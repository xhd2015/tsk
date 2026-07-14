# Scenario

**Feature**: top-level help lists channel subcommand

```
tsk --help -> channel in command list
```

## Steps

1. Run root help.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--help"}
	return nil
}
```
