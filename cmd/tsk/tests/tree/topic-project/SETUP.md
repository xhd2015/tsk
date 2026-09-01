# Scenario

**Feature**: under a topic, project is a secondary grouping level

```
create --topic eng "report"; set project.origin on task.json; tree --plain
```

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "report", "eng", nil)
	rel := readIndex(t, req, 1)
	taskDir := taskAbs(req, rel)
	task := readTaskJSON(t, taskDir)
	task.Project = &projectRefJSON{Origin: "github.com/xhd2015/wrk"}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(filepath.Join(taskDir, "task.json"), data, 0o644); err != nil {
		return err
	}
	req.Args = []string{"tree", "--plain"}
	return nil
}
```
