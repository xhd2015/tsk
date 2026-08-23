# Scenario

**Feature**: `tsk topic notes --json` emits a JSON array

```
note once -> notes --json
```

## Steps

1. mkdir + note `hello`.
2. `tsk topic notes --json knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "note", "knowledge-base", "hello")
	req.Args = []string{"topic", "notes", "--json", "knowledge-base"}
	return nil
}
```
