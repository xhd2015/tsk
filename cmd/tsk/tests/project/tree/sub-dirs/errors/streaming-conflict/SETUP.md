# Scenario

**Feature**: `--streaming` conflicts with `--no-streaming`

```
project tree --streaming --no-streaming -> Error, exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/example/root-repo.git")
	req.Args = []string{"project", "tree", "--streaming", "--no-streaming"}
	return nil
}
```
