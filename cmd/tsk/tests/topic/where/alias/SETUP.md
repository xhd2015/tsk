# Scenario

**Feature**: `tsk topic where` resolves an alias to the canonical directory

```
mkdir knowledge-base -> alias add 知识库 -> where 知识库
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic alias add knowledge-base 知识库`.
3. `tsk topic where 知识库`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "alias", "add", "knowledge-base", "知识库")
	req.Topic = "知识库"
	req.Args = []string{"topic", "where", "知识库"}
	return nil
}
```
