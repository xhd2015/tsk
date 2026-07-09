package tskcli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func runFollowup(home string, args []string) error {
	setCommand(currentCtx, "followup", args)

	if len(args) < 2 {
		return fail(fmt.Errorf("tsk followup: task id and message required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	message := args[1]
	if len(args) > 2 {
		message = joinArgs(args[1:])
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if task.Stage != "summary" {
		return fail(fmt.Errorf("invalid transition: followup only from summary"))
	}

	contextDir := filepath.Join(taskDir, "context")
	if err := os.MkdirAll(contextDir, 0o755); err != nil {
		return err
	}
	ts := storage.NowTimestamp(task.ID)
	filename := fmt.Sprintf("followup-%s.md", ts)
	content := message
	if content != "" && content[len(content)-1] != '\n' {
		content += "\n"
	}
	if err := os.WriteFile(filepath.Join(contextDir, filename), []byte(content), 0o644); err != nil {
		return err
	}

	_, err = storage.RenameTaskDir(home, &task, taskDir, "user_followup", "")
	return err
}