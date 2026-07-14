# Scenario

**Feature**: archived channel returns error

```
# index points to archive/<id> -> Error: ... exit 1
```

## Steps

1. Seed archived channel layout.
2. Run check.

```go
func Setup(t *testing.T, req *Request) error {
	writeArchivedChannel(t, req)
	req.LastActivity = oldActivityTS
	req.Args = defaultCheckArgs(req)
	return nil
}
```