# Scenario

**Feature**: `tsk topic view --json` is a nested object

```
mkdir -> view --json
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic view --json knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	req.Args = []string{"topic", "view", "--json", "knowledge-base"}
	return nil
}
```
