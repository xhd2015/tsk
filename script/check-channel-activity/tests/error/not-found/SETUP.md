# Scenario

**Feature**: missing channel returns error

```
# no channels/index entry -> Error: ... exit 1
```

## Steps

1. Do not seed channel fixtures.
2. Run check for `eng-alerts`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = defaultCheckArgs(req)
	return nil
}
```