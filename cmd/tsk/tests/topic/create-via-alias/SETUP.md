# Scenario

**Feature**: `tsk create --topic` stores the canonical path when given an alias

```
alias 知识库 -> create --topic 知识库 "x" -> topics/knowledge-base/…
```

## Steps

1. mkdir `knowledge-base`.
2. alias add `知识库`.
3. `tsk create --topic 知识库 "x"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "alias", "add", "knowledge-base", "知识库")
	req.Title = "x"
	req.Topic = "知识库"
	req.Args = []string{"create", "--topic", "知识库", "x"}
	return nil
}
```
