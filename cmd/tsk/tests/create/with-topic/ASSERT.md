## Expected

- Exit code 0.
- Task directory under `topics/eng/backend/1-create-x/`.
- `index/1` points at the topic-relative path.
- `task.json` has `topic_path: ["engineering","backend"]` or `["eng","backend"]` matching the topic segments.

## Side Effects

- Topic parent directories created under `topics/`.

## Exit Code

- 0

```go
import (
	"path/filepath"
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	wantRel := topicTaskRel(req.Topic, 1, "create", req.Title)
	taskDir := taskAbs(req, wantRel)
	assertDirExists(t, taskDir)
	assertFileExists(t, filepath.Join(taskDir, "task.json"))
	assertIndexEquals(t, req, 1, wantRel)

	task := readTaskJSON(t, taskDir)
	if task.Stage != "create" {
		t.Fatalf("stage: got %q want create", task.Stage)
	}
	wantTopic := strings.Split(req.Topic, "/")
	assertTopicPathEquals(t, req, 1, wantTopic)
}
```