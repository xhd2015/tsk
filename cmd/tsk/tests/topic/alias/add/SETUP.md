# Scenario

**Feature**: alias add creates topic.json with aliases

```
mkdir knowledge-base -> alias add 知识库 -> topic.json
```

## Steps

1. mkdir `knowledge-base`.
2. `tsk topic alias add knowledge-base 知识库`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	req.Args = []string{"topic", "alias", "add", "knowledge-base", "知识库"}
	return nil
}
```
