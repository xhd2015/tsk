# Scenario

**Feature**: `--streaming` conflicts with `--json`

```
project tree --streaming --json -> Error, exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/example/root-repo.git")
	req.Args = []string{"project", "tree", "--streaming", "--json"}
	return nil
}
```
