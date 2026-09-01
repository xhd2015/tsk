# Scenario

**Feature**: `tsk project add --note` seeds notes.jsonl

```
tsk project add --note "ctx" "titled" -> notes count 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	req.Title = "titled"
	req.Args = []string{"project", "add", "--note", "ctx", req.Title}
	return nil
}
```
