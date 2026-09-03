# Scenario

**Feature**: sub-dir flags conflict with `--all`

```
project tree --all --no-sub-dirs -> Error, exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/example/root-repo.git")
	req.Args = []string{"project", "tree", "--all", "--no-sub-dirs"}
	return nil
}
```
