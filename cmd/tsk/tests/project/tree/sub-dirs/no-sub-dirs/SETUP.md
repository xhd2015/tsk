# Scenario

**Feature**: `--no-sub-dirs` keeps current-project-only tree

```
setupNestedProjectTree; project tree --no-sub-dirs --plain -> only root-repo
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	setupNestedProjectTree(t, req)
	req.Args = []string{"project", "tree", "--no-sub-dirs", "--plain"}
	return nil
}
```
