# Scenario

**Feature**: `tsk topic info --json` is a machine-readable object

```
mkdir knowledge-base -> topic info --json knowledge-base
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic info --json knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	req.Args = []string{"topic", "info", "--json", "knowledge-base"}
	return nil
}
```
