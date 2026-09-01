# Scenario

**Feature**: `tsk list --topic` matches an alias to the canonical prefix

```
create --topic 知识库 -> list --topic 知识库 prints the id
```

## Steps

1. mkdir + alias `知识库`.
2. create `--topic 知识库 "x"`.
3. `tsk list --topic 知识库`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	runTskOK(t, req, "topic", "alias", "add", "knowledge-base", "知识库")
	id := addTask(t, req, "x", "知识库", nil)
	req.TaskID = id
	req.Args = []string{"list", "--topic", "知识库"}
	return nil
}
```
