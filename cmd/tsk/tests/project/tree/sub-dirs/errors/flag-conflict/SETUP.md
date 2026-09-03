# Scenario

**Feature**: `--no-sub-dirs` conflicts with `--sub-dirs-depth`

```
project tree --no-sub-dirs --sub-dirs-depth 2 -> Error, exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/example/root-repo.git")
	req.Args = []string{"project", "tree", "--no-sub-dirs", "--sub-dirs-depth", "2"}
	return nil
}
```
