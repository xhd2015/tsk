# Scenario

**Feature**: default (and `--streaming`) prints root first, then discovered projects

```
setupNestedProjectTree; project tree --streaming --plain
  -> root-repo before nested-repo
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	setupNestedProjectTree(t, req)
	req.Args = []string{"project", "tree", "--streaming", "--plain"}
	return nil
}
```
