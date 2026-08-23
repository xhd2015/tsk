# Scenario

**Feature**: view resolves an alias to the canonical tree

```
alias 知识库 -> topic view 知识库
```

## Steps

1. mkdir `knowledge-base`.
2. alias add `知识库`.
3. `tsk topic view 知识库`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "alias", "add", "knowledge-base", "知识库")
	req.Args = []string{"topic", "view", "知识库"}
	return nil
}
```
