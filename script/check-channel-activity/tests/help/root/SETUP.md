# Scenario

**Feature**: root help lists required and optional flags

```
check-channel-activity -h -> Usage + --channel-id + --exec-if-idle-1h
```

## Steps

1. Run `-h`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	return nil
}
```