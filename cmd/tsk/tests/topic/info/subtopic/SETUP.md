# Scenario

**Feature**: `tsk topic info` lists child topic directories

```
mkdir knowledge-base + knowledge-base/reports -> info lists reports
```

## Steps

1. mkdir `knowledge-base`.
2. mkdir `knowledge-base/reports`.
3. `tsk topic info knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "mkdir", "knowledge-base/reports")
	req.Args = []string{"topic", "info", "knowledge-base"}
	return nil
}
```
