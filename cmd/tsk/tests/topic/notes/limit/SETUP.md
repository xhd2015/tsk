# Scenario

**Feature**: `--limit 1` prints only the last note

```
note one; note two -> notes --limit 1 -> two only
```

## Steps

1. mkdir; note `one`; note `two`.
2. `tsk topic notes --limit 1 knowledge-base`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "note", "knowledge-base", "one")
	runTskOK(t, req, "topic", "note", "knowledge-base", "two")
	req.Args = []string{"topic", "notes", "--limit", "1", "knowledge-base"}
	return nil
}
```
