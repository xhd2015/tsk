# Scenario

**Feature**: `--sub-dirs-depth 2` includes depth-2 nested, excludes depth-4 deep

```
setupNestedProjectTree; project tree --sub-dirs-depth 2 --plain
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	setupNestedProjectTree(t, req)
	req.Args = []string{"project", "tree", "--sub-dirs-depth", "2", "--plain"}
	return nil
}
```
