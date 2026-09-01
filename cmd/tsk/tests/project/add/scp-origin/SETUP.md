# Scenario

**Feature**: SCP-style `gitlab@host:path.git` origin normalizes to host/path

```
origin gitlab@git.example.com:acme/loan-service/credit_backend/code-lens/tools/widget-cli.git
tsk project add "x" -> project.origin host/path
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "gitlab@git.example.com:acme/loan-service/credit_backend/code-lens/tools/widget-cli.git")
	req.Title = "x"
	req.Args = []string{"project", "add", req.Title}
	return nil
}
```
