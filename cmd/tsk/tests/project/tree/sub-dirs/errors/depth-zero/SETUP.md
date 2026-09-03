# Scenario

**Feature**: `--sub-dirs-depth 0` is invalid

```
project tree --sub-dirs-depth 0 -> Error, exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/example/root-repo.git")
	req.Args = []string{"project", "tree", "--sub-dirs-depth", "0"}
	return nil
}
```
