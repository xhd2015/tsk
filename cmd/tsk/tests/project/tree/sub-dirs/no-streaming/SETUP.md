# Scenario

**Feature**: `--no-streaming` buffers then prints (root first)

```
setupNestedProjectTree; project tree --no-streaming --plain
  -> root before nested; same membership
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	setupNestedProjectTree(t, req)
	req.Args = []string{"project", "tree", "--no-streaming", "--plain"}
	return nil
}
```
