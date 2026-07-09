# Scenario

**Feature**: agent format facts include slash-joined topic path above dir for topic-placed tasks

```
# create --topic stores topic_path segments; agent status prints topic: eng/backend above dir:
tsk create --topic eng/backend "topic status fact" -> tsk status --format=agent <id>
# facts: id → title → stage → terminal → topic: eng/backend → dir under topics/eng/backend/
```

## Steps

1. Create task with `--topic eng/backend` and a known title (stage `create`).
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "topic status fact"
	req.Topic = "eng/backend"
	id := createTask(t, req, req.Title, req.Topic, nil)
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
