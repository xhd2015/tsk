# Scenario

**Feature**: the same alias cannot belong to two topics

```
alias add a 知识库; alias add b 知识库 -> error
```

## Steps

1. mkdir `knowledge-base` and `other`.
2. alias add `knowledge-base` `知识库`.
3. `tsk topic alias add other 知识库`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "mkdir", "other")
	runTskOK(t, req, "topic", "alias", "add", "knowledge-base", "知识库")
	req.Args = []string{"topic", "alias", "add", "other", "知识库"}
	return nil
}
```
