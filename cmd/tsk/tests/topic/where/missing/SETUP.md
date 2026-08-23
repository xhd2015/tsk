# Scenario

**Feature**: `tsk topic where` errors when the topic does not exist

```
topic where no-such -> Error: topic not found
```

## Steps

1. `tsk topic where no-such` (no mkdir).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"topic", "where", "no-such"}
	return nil
}
```
