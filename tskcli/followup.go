package tskcli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runFollowup(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("followup", args)

	remaining, err := lessflags.
		Help("-h,--help", followupHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) < 2 {
		return fail(fmt.Errorf("tsk followup: task id and message required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}
	message := remaining[1]
	if len(remaining) > 2 {
		message = joinArgs(remaining[1:])
	}
	invk.setData(storage.EventData{TaskID: id, Text: message})

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

	invk.setData(storage.EventData{Stage: "user_followup"})
	return storage.SetTaskStage(&task, taskDir, "user_followup", "")
}
