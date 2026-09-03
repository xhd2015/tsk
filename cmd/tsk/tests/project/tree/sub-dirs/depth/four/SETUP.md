# Scenario

**Feature**: `--sub-dirs-depth 4` includes deep repo at depth 4

```
setupNestedProjectTree; project tree --sub-dirs-depth 4 --plain
  -> root + nested + deep (3 projects)
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	setupNestedProjectTree(t, req)
	req.Args = []string{"project", "tree", "--sub-dirs-depth", "4", "--plain"}
	return nil
}
```
