# Scenario

**Feature**: default `project tree` includes nested git projects (depth ≤ 3)

```
setupNestedProjectTree; project tree --plain
  -> root-repo + nested-repo; idle omitted; deep (depth 4) omitted
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	setupNestedProjectTree(t, req)
	req.Args = []string{"project", "tree", "--plain"}
	return nil
}
```
