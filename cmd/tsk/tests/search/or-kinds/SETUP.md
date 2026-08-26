# Scenario

**Feature**: `--task --note` ORs surfaces (title + note; not progress)

```
title has foo; note has bar; progress has baz -> search --task --note foo|bar hits; baz only via progress
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "alpha-foo-title", "", nil)
	runTskOK(t, req, "note", "add", "--id", fmt.Sprintf("%d", id), "note-bar-body")
	runTskOK(t, req, "progress", "add", "--id", fmt.Sprintf("%d", id), "--status", "in-progress", "progress-baz-body")
	req.Args = []string{"search", "--task", "--note", "bar"}
	return nil
}
```
