# Scenario

**Feature**: `--streaming` conflicts with `--all`

```
project tree --streaming --all -> Error, exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/example/root-repo.git")
	req.Args = []string{"project", "tree", "--streaming", "--all"}
	return nil
}
```
